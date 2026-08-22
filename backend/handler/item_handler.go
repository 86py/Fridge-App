package handler

import (
	"encoding/json"
	"fridge-backend/database"
	"fridge-backend/model"
	"net/http"
)

func HandleItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "GET" {
		rows, err := database.DB.Query("SELECT id, name, quantity, category, expiration_date, created_date FROM items ORDER BY expiration_date DESC")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []model.Item
		for rows.Next() {
			var item model.Item
			if err := rows.Scan(&item.ID, &item.Name, &item.Quantity, &item.Category, &item.ExpirationDate, &item.CreatedAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)

	} else if r.Method == "POST" {
		var item model.Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := database.DB.Exec(
			"INSERT INTO items (name, quantity, category, expiration_date) VALUES (?, ?, ?, ?)",
			item.Name, item.Quantity, item.Category, item.ExpirationDate,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "created successfully"})
	}
}
