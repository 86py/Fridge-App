# Fridge App

冷蔵庫の在庫管理アプリ。食材の登録・一覧表示・消費に加え、作った料理と使った食材を記録できる学習用プロジェクトです。

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
npm run test     # テスト実行 (Vitest)
npm run lint      # Lint (oxlint)
```

## 現状の実装範囲

- 実装済み:
  - 食材の登録（`POST /items`）、一覧取得（`GET /items`、賞味期限昇順）
  - 食材の消費（`POST /items/edit`）。消費後の数量が0以下ならレコード削除、在庫を超える消費は拒否
  - 賞味期限が近い/切れた食材のハイライト表示
  - 料理記録（`POST /dishes`・`GET /dishes`）。料理名・作った日・使った食材（複数）を1回の操作で登録し、対応する食材の在庫をまとめて消費する。料理カレンダーで日付ごとに振り返れる
- 未実装: 食材・料理記録の更新・削除、認証
- フロントエンドのAPI URLは`http://localhost:8080`にハードコード

## 注意点

- `backend/fridge.db`はSQLiteのDBファイル（起動時に自動生成）。git管理対象外
- CORSは全オリジン許可（`Access-Control-Allow-Origin: *`）— 開発用設定のため本番では要見直し
- コメント・UIラベルは日本語で統一
