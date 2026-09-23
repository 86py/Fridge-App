# CLAUDE.md

プロジェクト概要・技術構成・開発コマンドは[README.md](README.md)を参照。

## 構成（コード探索用）

- backend/main.go: エントリーポイント。:8080 でHTTPサーバー起動、ルーティング登録
- backend/database/db.go: SQLite接続・テーブル自動作成（items/dishes/dish_ingredients）
- backend/database/consume.go: 在庫消費ロジック（ConsumeItem）。複数箇所から呼べるよう独立させている
- backend/model/item.go, dish.go: Item/Dish/DishIngredientモデル定義
- backend/handler/item_handler.go: /items のGET（一覧）・POST（追加）ハンドラー
- backend/handler/item_edit_handler.go: /items/edit の在庫消費ハンドラー
- backend/handler/dish_handler.go: /dishes の料理記録作成・一覧ハンドラー
- frontend/src/App.jsx: 状態管理とAPI呼び出しを持つ親コンポーネント
- frontend/src/components/: ItemForm/ItemList/DishForm/DishCalendarに分割済み
- frontend/src/utils/expiration.js, calendar.js: 賞味期限判定・カレンダーグリッド生成の純粋関数

## 注意点

- backend/fridge.db はSQLiteのDBファイル（起動時に自動生成）。誤ってコミット・削除しないよう注意
- グローバル規約（~/.claude/CLAUDE.md）に準拠：新規コードは理由をdocstringに記載し、TDDで進めること
  - バックエンド・フロントエンドともにテストを整備済み（`go test ./...`, `npm run test`）
