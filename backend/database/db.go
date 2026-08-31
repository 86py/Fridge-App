package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// CreateTables はitemsテーブルを作成する。本番用DB・テスト用インメモリDBの
// 両方から呼び出し、スキーマ定義を1箇所に集約するために切り出している。
func CreateTables(db *sql.DB) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		quantity INTEGER,
		category TEXT,
		expiration_date TEXT,
		created_date DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTableSQL)
	return err
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
