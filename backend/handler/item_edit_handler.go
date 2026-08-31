package handler

import (
	"encoding/json"
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

		// 1. 現在の数量を確認する
		var quantity int
		err := database.DB.QueryRow("SELECT quantity FROM items WHERE id = ?", req.ID).Scan(&quantity)
		if err != nil {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		// 2. 在庫を超える消費はエラーにする（押せなくする制御はフロント側の責務とし、
		// サーバー側は不正なリクエストの最終防衛ラインとしてのみチェックする）
		if req.Quantity > quantity {
			http.Error(w, "quantity exceeds current stock", http.StatusBadRequest)
			return
		}

		// 3. 減算後の数量で分岐（0になるなら削除、そうでなければ更新）
		newQuantity := quantity - req.Quantity
		if newQuantity > 0 {
			_, err = database.DB.Exec("UPDATE items SET quantity = ? WHERE id = ?", newQuantity, req.ID)
		} else {
			_, err = database.DB.Exec("DELETE FROM items WHERE id = ?", req.ID)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "updated successfully"})
	}
}
