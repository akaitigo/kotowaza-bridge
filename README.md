# Kotowaza Bridge

[![CI](https://github.com/akaitigo/kotowaza-bridge/actions/workflows/ci.yml/badge.svg)](https://github.com/akaitigo/kotowaza-bridge/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

日本のことわざと世界各国の類似表現を対比して学ぶ多言語対照辞典アプリ。
「猿も木から落ちる」と英語 "Even Homer sometimes nods" のように、同じ教訓を各文化がどんな比喩で表現するかを探索できます。

## 機能

- **ことわざ対照表示** — 日本語ことわざ20選と英・中・韓の対応表現を並べて表示
- **LLMロールプレイ** — ことわざの自然な使い方をAIとの会話で練習
- **文化背景解説** — 各ことわざの由来・使用例・文化的背景を解説
- **キーワード検索** — リアルタイム検索でことわざを素早く見つける

## Tech Stack

| レイヤー | 技術 |
|---------|------|
| Frontend | TypeScript / Next.js 15 (App Router) |
| API | Go / chi router |
| Database | PostgreSQL |
| LLM | Claude API (Anthropic) |
| CI | GitHub Actions |

## クイックスタート

### 前提条件

- Go 1.23+
- Node.js 22+
- PostgreSQL 15+

### セットアップ

```bash
# リポジトリをクローン
git clone git@github.com:akaitigo/kotowaza-bridge.git
cd kotowaza-bridge

# 環境変数を設定
cp .env.example .env
# .env を編集して DATABASE_URL と LLM_API_KEY を設定

# データベースをセットアップ
# 注意: 002 はことわざ本体を参照する追加データのため、seed の後に適用する
psql $DATABASE_URL < api/db/migrations/001_create_tables.sql
psql $DATABASE_URL < api/db/seed.sql
psql $DATABASE_URL < api/db/migrations/002_additional_equivalents.sql
psql $DATABASE_URL < api/db/migrations/003_search_indexes.up.sql

# API サーバーを起動
cd api && go run ./cmd/server

# フロントエンドを起動（別ターミナル）
cd web && npm install && npm run dev
```

http://localhost:3000 でアプリにアクセスできます。

## アーキテクチャ

```
kotowaza-bridge/
├── api/                    # Go API サーバー
│   ├── cmd/server/         # エントリーポイント
│   ├── internal/
│   │   ├── config/         # 環境変数ベースの設定
│   │   ├── domain/         # ドメインモデル
│   │   ├── handler/        # HTTPハンドラー
│   │   ├── middleware/     # CORS等のミドルウェア
│   │   ├── repository/     # データアクセス層
│   │   └── service/        # ビジネスロジック + LLMクライアント
│   └── db/                 # マイグレーション + シードデータ
├── web/                    # Next.js フロントエンド
│   └── src/
│       ├── app/            # App Router ページ
│       ├── components/     # UIコンポーネント
│       ├── lib/            # API クライアント
│       └── types/          # TypeScript 型定義
└── docs/                   # ADR、ドキュメント
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/health` | Liveness（プロセス稼働確認、依存先に触れない） |
| GET | `/api/v1/health/ready` | Readiness（DB接続を検証。障害時は503） |
| GET | `/api/v1/kotowaza` | ことわざ一覧 |
| GET | `/api/v1/kotowaza/:id` | ことわざ詳細 |
| GET | `/api/v1/kotowaza/search?q=` | キーワード検索 |
| POST | `/api/v1/chat` | LLMチャット |

## ライセンス

MIT
