# Changelog

## [1.0.0] - 2026-04-04

### Added
- ことわざ対照辞典の初期リリース
- Go API サーバー (chi router + pgx)
  - ことわざ一覧・詳細・検索エンドポイント
  - LLMチャットエンドポイント (Anthropic Claude API)
  - ヘルスチェックエンドポイント
- Next.js フロントエンド (App Router)
  - ことわざ一覧ページ（カード表示 + リアルタイム検索）
  - ことわざ詳細ページ（対応表現・文化背景・由来）
  - LLMロールプレイチャット画面
- PostgreSQL スキーマ + シードデータ
  - 日本語ことわざ20件
  - 英語・中国語・韓国語の対応表現60件
- CI/CD (GitHub Actions)
- ADR-001: LLM統合設計方針
