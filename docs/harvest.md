# Harvest: Kotowaza Bridge

## プロジェクト概要
- **アイデアID**: #1812
- **名前**: Kotowaza-Bridge
- **ドメイン**: ことわざ学習
- **完了日**: 2026-04-04

## メトリクス

| 項目 | 値 |
|------|-----|
| コミット数 | 9 |
| PR数 | 8 (全MERGED) |
| Issue数 | 5 (全closed) |
| ADR数 | 1 |
| CLAUDE.md行数 | 34 (上限50以下) |
| Go ソースファイル | 14 |
| Go テストファイル | 5 |
| TypeScript ファイル | 9 |

## Review Loop 履歴

| ラウンド | 焦点 | CRITICAL | HIGH | MEDIUM | LOW | 修正PR |
|---------|------|----------|------|--------|-----|--------|
| R1 | 表面: lint, naming, imports | 0 | 4 | 5 | 3 | #12 |
| R2-R5 | 実装+設計+統合+堅牢性 | 0 | 4 | 7 | 1 | #13 |
| R6 | 副作用確認 | 0 | 0 | 0 | 0 | - |

### R1 主要修正
- `interface{}` → `any` (Go 1.18+ convention)
- `errors.Is()` for error comparison
- `io.LimitReader` for LLM response size cap
- `http.MaxBytesReader` for request body size limit
- Chat message validation (role whitelist, length/count limits)
- Unused import removal, API client consolidation

### R2-R5 主要修正
- Sentinel errors (`ErrValidation`, `ErrNotFound`) for HTTP status differentiation
- Chat handler: 400/404/500 proper status codes
- Validation test coverage: 4 new test cases
- Graceful shutdown: `errCh` pattern replacing `log.Fatalf` in goroutine

## 技術的判断

### 良かった点
- Go の clean architecture (domain/service/repository) が明確に分離
- interface 抽象化で LLM クライアントのモック注入が容易
- Next.js App Router + RSC のシンプルな構成
- パラメータバインドで SQL injection を完全防止

### 改善提案
- E2E テストの追加（Playwright）
- LLM レスポンスのストリーミング対応
- PostgreSQL tsvector を使った全文検索（現在はLIKE）
- Rate limiting ミドルウェアの実装

## テンプレート改善提案

1. **Go post-lint hook**: `errcheck` がCIで落ちる前にローカルで検出するために、golangci-lint の設定にerrcheckを明示的に含めるべき
2. **biome organizeImports**: 初回コミット時にformat + checkを実行するステップをLaunch手順に追加
3. **Review Loop R1-R6**: R2-R5を1回のエージェントに統合して効率化。R1（表面）→ R2-R5（深層統合）→ R6（副作用確認）の3段階が実用的
