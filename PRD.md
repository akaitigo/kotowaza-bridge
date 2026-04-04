# PRD: Kotowaza-Bridge

## 概要
日本のことわざと世界各国の類似表現を対比して学ぶ多言語対照辞典アプリ。
「猿も木から落ちる」と英語 "Even Homer sometimes nods"、韓国語・中国語の対応表現を並べて文化比較できる。

## 背景
ことわざは各文化の価値観や世界観を凝縮した表現。同じ教訓でも国ごとに異なる比喩を用いることが多い。
これを対比的に学ぶことで、言語学習と異文化理解を同時に深められる。

## ユニークポイント
同義のことわざを世界30カ国から収集し、「同じ教訓を各文化がどんな比喩で表現するか」を比較文化的に探索できるクロスカルチャー辞典。

## Tech Stack
- **Frontend**: TypeScript / Next.js (App Router)
- **API**: Go (net/http + chi router)
- **Database**: PostgreSQL
- **Search**: 全文検索（PostgreSQL tsvector）
- **LLM**: Claude API（ロールプレイ練習用）

## MVP機能

### F1: ことわざ対照表示
- [x] 日本語ことわざ一覧表示（300選のうちMVPは初期シード20件）
- [x] 各ことわざに対する英語・中国語・韓国語の対応表現を並べて表示
- [x] ことわざの意味・由来・使用例の表示
- [x] キーワード検索

### F2: LLMロールプレイ練習
- [x] ことわざを選んでLLMとの会話練習を開始
- [x] LLMが日常会話シナリオを生成し、ことわざの自然な使い方を練習
- [x] 会話履歴の表示

### F3: 文化背景解説
- [x] 各ことわざの文化的背景の解説表示
- [x] 類似表現の言語別マッピング表示

## 非機能要件
- API レスポンス: p95 < 200ms（DB クエリ）
- LLM レスポンス: ストリーミング対応
- モバイルレスポンシブ対応

## API設計

### Endpoints
- `GET /api/v1/kotowaza` — ことわざ一覧取得（ページネーション対応）
- `GET /api/v1/kotowaza/:id` — ことわざ詳細取得（対応表現含む）
- `GET /api/v1/kotowaza/search?q=` — キーワード検索
- `POST /api/v1/chat` — LLMロールプレイチャット
- `GET /api/v1/health` — ヘルスチェック

## DB設計

### kotowaza テーブル
| Column | Type | Description |
|--------|------|-------------|
| id | UUID | PK |
| japanese | TEXT | 日本語ことわざ |
| reading | TEXT | 読み仮名 |
| meaning | TEXT | 意味 |
| origin | TEXT | 由来 |
| usage_example | TEXT | 使用例 |
| cultural_note | TEXT | 文化的背景 |
| created_at | TIMESTAMPTZ | 作成日時 |

### equivalent テーブル
| Column | Type | Description |
|--------|------|-------------|
| id | UUID | PK |
| kotowaza_id | UUID | FK → kotowaza |
| language | TEXT | 言語コード (en, zh, ko) |
| expression | TEXT | 対応表現 |
| literal_meaning | TEXT | 直訳 |
| explanation | TEXT | 解説 |
