# CLAUDE.md

プロジェクト概要・技術構成・開発コマンドは[README.md](README.md)を参照。

## 構成（コード探索用）

- backend/main.go: エントリーポイント。:8080 でHTTPサーバー起動
- backend/database/db.go: SQLite接続・テーブル自動作成（items）
- backend/model/item.go: Itemモデル定義
- backend/handler/item_handler.go: /items のGET（一覧）・POST（追加）ハンドラー
- frontend/src/App.jsx: 単一コンポーネントでフォーム入力・一覧表示を実装

## 注意点

- backend/fridge.db はSQLiteのDBファイル（起動時に自動生成）。誤ってコミット・削除しないよう注意
- グローバル規約（~/.claude/CLAUDE.md）に準拠：新規コードは理由をdocstringに記載し、TDDで進めること
  - 現状はdocstring未記載・テストカバレッジ不足（OPTIONSのみ）のため、今後の変更で順次追従する
