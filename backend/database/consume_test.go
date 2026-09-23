package database

import (
	"database/sql"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func newConsumeTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	if err := CreateTables(db); err != nil {
		t.Fatal(err)
	}
	return db
}

func insertConsumeTestItem(t *testing.T, db *sql.DB, quantity int) int {
	t.Helper()
	res, err := db.Exec(
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
	return int(id)
}

func TestConsumeItem_UpdatesQuantityWhenStockRemains(t *testing.T) {
	db := newConsumeTestDB(t)
	id := insertConsumeTestItem(t, db, 5)

	tx, _ := db.Begin()
	if err := ConsumeItem(tx, id, 3); err != nil {
		t.Fatalf("想定外のエラー: %v", err)
	}
	tx.Commit()

	var quantity int
	if err := db.QueryRow("SELECT quantity FROM items WHERE id = ?", id).Scan(&quantity); err != nil {
		t.Fatalf("レコードが見つかりません: %v", err)
	}
	if quantity != 2 {
		t.Errorf("数量が期待値と異なります: got %v want %v", quantity, 2)
	}
}

func TestConsumeItem_DeletesItemWhenFullyConsumed(t *testing.T) {
	db := newConsumeTestDB(t)
	id := insertConsumeTestItem(t, db, 5)

	tx, _ := db.Begin()
	if err := ConsumeItem(tx, id, 5); err != nil {
		t.Fatalf("想定外のエラー: %v", err)
	}
	tx.Commit()

	var quantity int
	err := db.QueryRow("SELECT quantity FROM items WHERE id = ?", id).Scan(&quantity)
	if err == nil {
		t.Errorf("在庫を使い切ったレコードが削除されていません: quantity=%v", quantity)
	}
}

func TestConsumeItem_ReturnsErrInsufficientStock(t *testing.T) {
	db := newConsumeTestDB(t)
	id := insertConsumeTestItem(t, db, 3)

	tx, _ := db.Begin()
	err := ConsumeItem(tx, id, 10)
	tx.Rollback()

	if !errors.Is(err, ErrInsufficientStock) {
		t.Errorf("エラーが期待値と異なります: got %v want %v", err, ErrInsufficientStock)
	}
}

func TestConsumeItem_ReturnsErrItemNotFound(t *testing.T) {
	db := newConsumeTestDB(t)

	tx, _ := db.Begin()
	err := ConsumeItem(tx, 9999, 1)
	tx.Rollback()

	if !errors.Is(err, ErrItemNotFound) {
		t.Errorf("エラーが期待値と異なります: got %v want %v", err, ErrItemNotFound)
	}
}

func TestConsumeItem_RollbackLeavesStockUnchanged(t *testing.T) {
	db := newConsumeTestDB(t)
	id := insertConsumeTestItem(t, db, 3)

	tx, _ := db.Begin()
	_ = ConsumeItem(tx, id, 10) // 在庫不足で失敗する想定
	tx.Rollback()

	var quantity int
	if err := db.QueryRow("SELECT quantity FROM items WHERE id = ?", id).Scan(&quantity); err != nil {
		t.Fatalf("レコードが見つかりません: %v", err)
	}
	if quantity != 3 {
		t.Errorf("ロールバック後に数量が変化しています: got %v want %v", quantity, 3)
	}
}
