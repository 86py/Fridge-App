# CLAUDE.md

## プロジェクト概要
冷蔵庫の在庫管理アプリ。食材の登録・一覧表示を行う。
- backend/: Go製REST API（SQLite使用）
- frontend/: React + Vite製SPA

## 構成
- backend/main.go: エントリーポイント。:8080 でHTTPサーバー起動
- backend/database/db.go: SQLite接続・テーブル自動作成（items）
- backend/model/item.go: Itemモデル定義
- backend/handler/item_handler.go: /items のGET（一覧）・POST（追加）ハンドラー
- frontend/src/App.jsx: 単一コンポーネントでフォーム入力・一覧表示を実装

## 現状の実装範囲
- 実装済み: 食材の登録（POST /items）、一覧取得（GET /items、賞味期限降順）
- 未実装: 更新・削除（PUT/DELETE）、入力バリデーション、認証
- フロントエンドのAPI URLは http://localhost:8080 にハードコード

## 開発コマンド
### backend
- 起動: `go run main.go` (backend/ ディレクトリ内)
- テスト: `go test ./...`

### frontend
- 開発サーバー: `npm run dev`
- ビルド: `npm run build`
- Lint: `npm run lint` (oxlint)

## 注意点
- backend/fridge.db はSQLiteのDBファイル（起動時に自動生成）。誤ってコミット・削除しないよう注意
- CORSは全オリジン許可（Access-Control-Allow-Origin: *）— 開発用設定のため本番では要見直し
- コメントやUIラベルは日本語で統一されている
- グローバル規約（~/.claude/CLAUDE.md）に準拠：新規コードは理由をdocstringに記載し、TDDで進めること
  - 現状はdocstring未記載・テストカバレッジ不足（OPTIONSのみ）のため、今後の変更で順次追従する
