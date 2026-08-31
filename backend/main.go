package main

import (
	"fridge-backend/database"
	"fridge-backend/handler"
	"log"
	"net/http"
)

func main() {
	// データベースの初期化
	database.InitDB()
	defer database.DB.Close()

	// ルーティング設定
	http.HandleFunc("/items", handler.HandleItems)
	http.HandleFunc("/items/edit", handler.HandleItemsEdit)
	log.Println("Backend server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
