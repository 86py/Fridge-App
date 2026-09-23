package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fridge-backend/database"
	"fridge-backend/model"
	"net/http"
)

func HandleDishes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "GET" {
		handleGetDishes(w, r)
	} else if r.Method == "POST" {
		handleCreateDish(w, r)
	}
}

func handleGetDishes(w http.ResponseWriter, r *http.Request) {
	// dishes本体を先に全件読み切ってrowsをCloseしてから、各dishのingredientsを
	// 個別にクエリする。rowsをオープンしたまま同じDBハンドルで別クエリを発行すると、
	// コネクションプールの空きを待って(接続数の設定によっては)デッドロックしうるため。
	rows, err := database.DB.Query("SELECT id, name, cooked_date, created_date FROM dishes ORDER BY cooked_date DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dishes := []model.Dish{}
	for rows.Next() {
		var dish model.Dish
		if err := rows.Scan(&dish.ID, &dish.Name, &dish.CookedDate, &dish.CreatedAt); err != nil {
			rows.Close()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dishes = append(dishes, dish)
	}
	rows.Close()

	for i := range dishes {
		ingredientRows, err := database.DB.Query("SELECT item_name, quantity FROM dish_ingredients WHERE dish_id = ?", dishes[i].ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		dishes[i].Ingredients = []model.DishIngredient{}
		for ingredientRows.Next() {
			var ingredient model.DishIngredient
			if err := ingredientRows.Scan(&ingredient.ItemName, &ingredient.Quantity); err != nil {
				ingredientRows.Close()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dishes[i].Ingredients = append(dishes[i].Ingredients, ingredient)
		}
		ingredientRows.Close()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dishes)
}

func handleCreateDish(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		CookedDate  string `json:"cooked_date"`
		Ingredients []struct {
			ItemID   int `json:"item_id"`
			Quantity int `json:"quantity"`
		} `json:"ingredients"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.CookedDate == "" || len(req.Ingredients) == 0 {
		http.Error(w, "name, cooked_date and ingredients are required", http.StatusBadRequest)
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 各食材について、消費前に名前をスナップショットとして取得しつつ、
	// ConsumeItemで在庫を減らす。1つでも失敗したら全体をロールバックし、
	// 一部の食材だけ消費されてしまう不整合を防ぐ。
	type consumedIngredient struct {
		Name     string
		Quantity int
	}
	consumed := make([]consumedIngredient, 0, len(req.Ingredients))

	for _, ing := range req.Ingredients {
		var itemName string
		err := tx.QueryRow("SELECT name FROM items WHERE id = ?", ing.ItemID).Scan(&itemName)
		if err == sql.ErrNoRows {
			tx.Rollback()
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := database.ConsumeItem(tx, ing.ItemID, ing.Quantity); err != nil {
			tx.Rollback()
			switch {
			case errors.Is(err, database.ErrItemNotFound):
				http.Error(w, "Item not found", http.StatusNotFound)
			case errors.Is(err, database.ErrInsufficientStock):
				http.Error(w, "quantity exceeds current stock", http.StatusBadRequest)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		consumed = append(consumed, consumedIngredient{Name: itemName, Quantity: ing.Quantity})
	}

	result, err := tx.Exec("INSERT INTO dishes (name, cooked_date) VALUES (?, ?)", req.Name, req.CookedDate)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dishID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, c := range consumed {
		if _, err := tx.Exec(
			"INSERT INTO dish_ingredients (dish_id, item_name, quantity) VALUES (?, ?, ?)",
			dishID, c.Name, c.Quantity,
		); err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "created successfully", "id": dishID})
}
