package database

import (
	"database/sql"
	"errors"
)

// ErrItemNotFound は指定IDの食材が存在しない場合に返す。
var ErrItemNotFound = errors.New("item not found")

// ErrInsufficientStock は消費数が在庫数を超える場合に返す。
var ErrInsufficientStock = errors.New("quantity exceeds current stock")

// ConsumeItem はトランザクション内で指定IDの在庫をquantityだけ消費する。
// 減算後の数量が0以下になった場合はレコードごと削除し、そうでなければ更新する。
// 料理記録機能で複数食材を1トランザクションでまとめて消費できるよう、
// item_edit_handler.goにあった処理をハンドラーから独立させた。
func ConsumeItem(tx *sql.Tx, id int, quantity int) error {
	var currentQuantity int
	err := tx.QueryRow("SELECT quantity FROM items WHERE id = ?", id).Scan(&currentQuantity)
	if err == sql.ErrNoRows {
		return ErrItemNotFound
	}
	if err != nil {
		return err
	}

	if quantity > currentQuantity {
		return ErrInsufficientStock
	}

	newQuantity := currentQuantity - quantity
	if newQuantity > 0 {
		_, err = tx.Exec("UPDATE items SET quantity = ? WHERE id = ?", newQuantity, id)
	} else {
		_, err = tx.Exec("DELETE FROM items WHERE id = ?", id)
	}
	return err
}
