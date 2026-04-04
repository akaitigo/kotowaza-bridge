# kotowaza-bridge

## コマンド
- ビルド: `make build`
- テスト: `make test`
- lint: `make lint`
- フォーマット: `make format`
- 全チェック: `make check`

## ワークフロー
1. research.md → plan.md → 承認後に実装
2. plan.md のtodoで進捗管理

## 構造
- `api/` — Go API サーバー (chi router)
- `web/` — Next.js フロントエンド (App Router)
- `docs/adr/` — Architecture Decision Records

## ルール
- ADR: docs/adr/ 参照。新規決定はADRを書いてから実装
- テスト: 機能追加時は必ずテストを同時に書く
- lint設定の変更禁止（ADR必須）

## 禁止事項
- any型(TS) / console.log / TODO コメント
- .env・credentials のコミット
- lint設定の無効化

## Hooks
- 設定: .claude/hooks/ 参照

## 状態管理
- git log + GitHub Issues でセッション間の状態を管理
- セッション開始: `bash startup.sh`
