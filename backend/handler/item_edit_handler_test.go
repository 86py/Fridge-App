package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fridge-backend/database"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// テストごとにインメモリDBを用意し、本番用DBファイル(fridge.db)を汚さないようにする
func newTestDB(t *testing.T) {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	createTableSQL := `CREATE TABLE items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		quantity INTEGER,
		category TEXT,
		expiration_date TEXT,
		created_date DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		t.Fatal(err)
	}

	database.DB = db
}

func insertTestItem(t *testing.T, quantity int) int64 {
	t.Helper()
	res, err := database.DB.Exec(
		"INSERT INTO items (name, quantity, category, expiration_date) VALUES (?, ?, ?, ?)",
		"テスト食材", quantity, "野菜", "2026-12-31",
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func postDecrementRequest(t *testing.T, id int64, quantity int) *httptest.ResponseRecorder {
	t.Helper()
	body, _ := json.Marshal(map[string]int{"id": int(id), "quantity": quantity})
	req, err := http.NewRequest("POST", "/items/edit", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	HandleItemsEdit(rr, req)
	return rr
}

func TestHandleItemsEdit_UpdatesQuantityWhenStockRemains(t *testing.T) {
	newTestDB(t)
	id := insertTestItem(t, 5)

	rr := postDecrementRequest(t, id, 3)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("ステータスコードが期待値と異なります: got %v want %v", status, http.StatusOK)
	}

	var quantity int
	err := database.DB.QueryRow("SELECT quantity FROM items WHERE id = ?", id).Scan(&quantity)
	if err != nil {
		t.Fatalf("レコードが見つかりません: %v", err)
	}
	if quantity != 2 {
		t.Errorf("数量が期待値と異なります: got %v want %v", quantity, 2)
	}
}

func TestHandleItemsEdit_DeletesItemWhenFullyConsumed(t *testing.T) {
	newTestDB(t)
	id := insertTestItem(t, 5)

	rr := postDecrementRequest(t, id, 5)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("ステータスコードが期待値と異なります: got %v want %v", status, http.StatusOK)
	}

	var quantity int
	err := database.DB.QueryRow("SELECT quantity FROM items WHERE id = ?", id).Scan(&quantity)
	if err == nil {
		t.Errorf("在庫を使い切ったレコードが削除されていません: quantity=%v", quantity)
	}
}

func TestHandleItemsEdit_RejectsQuantityExceedingStock(t *testing.T) {
	newTestDB(t)
	id := insertTestItem(t, 3)

	rr := postDecrementRequest(t, id, 10)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("ステータスコードが期待値と異なります: got %v want %v", status, http.StatusBadRequest)
	}

	var quantity int
	err := database.DB.QueryRow("SELECT quantity FROM items WHERE id = ?", id).Scan(&quantity)
	if err != nil {
		t.Fatalf("レコードが見つかりません: %v", err)
	}
	if quantity != 3 {
		t.Errorf("在庫を超える消費を拒否した際に数量が変わっています: got %v want %v", quantity, 3)
	}
}

func TestHandleItemsEdit_RejectsNonPositiveQuantity(t *testing.T) {
	newTestDB(t)
	id := insertTestItem(t, 5)

	rr := postDecrementRequest(t, id, 0)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("ステータスコードが期待値と異なります: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestHandleItemsEdit_ReturnsNotFoundForMissingItem(t *testing.T) {
	newTestDB(t)

	rr := postDecrementRequest(t, 9999, 1)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("ステータスコードが期待値と異なります: got %v want %v", status, http.StatusNotFound)
	}
}
