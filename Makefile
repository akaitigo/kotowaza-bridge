.PHONY: build test lint format check clean quality harvest

# === API (Go) ===
build-api:
	cd api && go build -trimpath -ldflags "-s -w" ./cmd/server

test-api:
	cd api && go test -v -race -count=1 -coverprofile=coverage.out ./...

lint-api:
	cd api && golangci-lint run ./...

format-api:
	cd api && gofumpt -w . && goimports -w .

# === Web (Next.js) ===
build-web:
	cd web && npx next build

test-web:
	cd web && npx vitest run --passWithNoTests

lint-web:
	cd web && npx oxlint . && npx biome check .

format-web:
	cd web && npx biome format --write .

# === Combined ===
build: build-api build-web

test: test-api test-web

lint: lint-api lint-web

format: format-api format-web

check: format lint test build
	@echo "All checks passed."

quality:
	@echo "=== Quality Gate ==="
	@test -f LICENSE || { echo "ERROR: LICENSE missing. Fix: add MIT LICENSE file"; exit 1; }
	@! grep -rn "TODO\|FIXME\|HACK\|console\.log\|println\|print(" api/internal/ web/src/ 2>/dev/null | grep -v "node_modules" || { echo "ERROR: debug output or TODO found. Fix: remove before ship"; exit 1; }
	@! grep -rn "password=\|secret=\|api_key=\|sk-\|ghp_" api/ web/src/ 2>/dev/null | grep -v '\$$' | grep -v "node_modules" || { echo "ERROR: hardcoded secrets. Fix: use env vars with no default"; exit 1; }
	@test ! -f CLAUDE.md || [ $$(wc -l < CLAUDE.md) -le 50 ] || { echo "ERROR: CLAUDE.md is $$(wc -l < CLAUDE.md) lines (max 50). Fix: remove build details, use pointers only"; exit 1; }
	@echo "OK: automated quality checks passed"

clean:
	cd api && go clean -cache -testcache && rm -f coverage.out
	cd web && rm -rf .next/ node_modules/.cache/

harvest:
	@echo "=== Harvest ==="
	@mkdir -p docs
	@echo "# Harvest: kotowaza-bridge" > docs/harvest.md
	@echo "" >> docs/harvest.md
	@echo "## メトリクス" >> docs/harvest.md
	@echo "| 項目 | 値 |" >> docs/harvest.md
	@echo "|------|-----|" >> docs/harvest.md
	@echo "| コミット数 | $$(git log --oneline --no-merges | wc -l) |" >> docs/harvest.md
	@echo "| ADR数 | $$(ls docs/adr/*.md 2>/dev/null | wc -l) |" >> docs/harvest.md
	@echo "| CLAUDE.md行数 | $$(wc -l < CLAUDE.md 2>/dev/null || echo 0) |" >> docs/harvest.md
	@echo "Harvest report generated: docs/harvest.md"
