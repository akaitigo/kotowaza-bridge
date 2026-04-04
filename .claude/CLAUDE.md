# プロジェクト固有コンテキスト

## API (Go)
- Router: chi v5
- DB: pgx v5 (PostgreSQL)
- 構造: cmd/server → internal/{handler,service,repository,domain}
- テスト: testing + testify
- エラー: カスタムAppError型でHTTPステータスをマッピング

## Web (Next.js)
- App Router (RSC + Client Components)
- スタイル: Tailwind CSS
- 状態管理: React hooks + Server Actions
- API通信: fetch (Go APIへプロキシ)

## DB
- マイグレーション: api/db/migrations/ に連番SQLファイル
- シード: api/db/seed.sql
