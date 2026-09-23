package handler

import (
	"encoding/json"
	"errors"
	"fridge-backend/database"
	"net/http"
)

func HandleItemsEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" { // または PATCH
		// リクエストから対象のIDを受け取る構造体などを用意
		var req struct {
			ID       int `json:"id"`
			Quantity int `json:"quantity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Quantity <= 0 {
			http.Error(w, "quantity must be positive", http.StatusBadRequest)
			return
		}

		tx, err := database.DB.Begin()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 押せなくする制御はフロント側の責務とし、サーバー側(ConsumeItem)は
		// 不正なリクエストの最終防衛ラインとしてのみチェックする
		if err := database.ConsumeItem(tx, req.ID, req.Quantity); err != nil {
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

		if err := tx.Commit(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "updated successfully"})
	}
}
