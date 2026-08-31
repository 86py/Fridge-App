# Fridge App

冷蔵庫の在庫管理アプリ。食材の登録・一覧表示を行う学習用プロジェクトです。

- `backend/`: Go製REST API（SQLite使用）
- `frontend/`: React + Vite製SPA

## 技術構成

- バックエンド: Go（標準の`net/http`）+ SQLite
- フロントエンド: React + Vite

## 開発コマンド

### backend

```sh
cd backend
go run main.go   # :8080 でHTTPサーバー起動
go test ./...    # テスト実行
```

### frontend

```sh
cd frontend
npm run dev      # 開発サーバー起動
npm run build    # ビルド
npm run lint      # Lint (oxlint)
```

## 現状の実装範囲

- 実装済み: 食材の登録（`POST /items`）、一覧取得（`GET /items`、賞味期限降順）
- 未実装: 更新・削除（PUT/DELETE）、入力バリデーション、認証
- フロントエンドのAPI URLは`http://localhost:8080`にハードコード

## 注意点

- `backend/fridge.db`はSQLiteのDBファイル（起動時に自動生成）。git管理対象外
- CORSは全オリジン許可（`Access-Control-Allow-Origin: *`）— 開発用設定のため本番では要見直し
- コメント・UIラベルは日本語で統一
