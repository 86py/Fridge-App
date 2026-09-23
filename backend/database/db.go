package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// CreateTables はitems/dishes/dish_ingredientsテーブルを作成する。本番用DB・
// テスト用インメモリDBの両方から呼び出し、スキーマ定義を1箇所に集約するために切り出している。
func CreateTables(db *sql.DB) error {
	createItemsTableSQL := `CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		quantity INTEGER,
		category TEXT,
		expiration_date TEXT,
		created_date DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createItemsTableSQL); err != nil {
		return err
	}

	// dishesは「作った料理」の記録。同じ料理名を複数回作ることを許容するため
	// nameにUNIQUE制約は付けない(日によって使う食材・量が変わりうるため)。
	createDishesTableSQL := `CREATE TABLE IF NOT EXISTS dishes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		cooked_date TEXT NOT NULL,
		created_date DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createDishesTableSQL); err != nil {
		return err
	}

	// dish_ingredientsは1つのdishに紐づく「使った食材の内訳」。item_idではなく
	// item_nameのスナップショットを持つ。消費によりitemsレコードが削除され得るため、
	// IDだけ持つと参照が宙に浮いてしまう。
	createDishIngredientsTableSQL := `CREATE TABLE IF NOT EXISTS dish_ingredients (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		dish_id INTEGER NOT NULL,
		item_name TEXT NOT NULL,
		quantity INTEGER NOT NULL
	);`
	if _, err := db.Exec(createDishIngredientsTableSQL); err != nil {
		return err
	}

	return nil
}

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./fridge.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := CreateTables(DB); err != nil {
		log.Fatal(err)
	}
}
