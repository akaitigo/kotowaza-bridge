#!/usr/bin/env bash
set -euo pipefail

echo "=== Session Startup ==="

[ -d ".git" ] || { echo "ERROR: Not in git repository"; exit 1; }

echo "=== Tool auto-install ==="
if [ -f "api/go.mod" ]; then
  echo "Detected: Go (api/)"
  command -v golangci-lint &>/dev/null || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest 2>/dev/null || echo "WARN: golangci-lint install failed"; }
  command -v gofumpt &>/dev/null || { echo "Installing gofumpt..."; go install mvdan.cc/gofumpt@latest 2>/dev/null || echo "WARN: gofumpt install failed"; }
fi
if [ -f "web/package.json" ]; then
  echo "Detected: TypeScript (web/)"
  command -v oxlint &>/dev/null || { echo "Installing oxlint..."; npm install -g oxlint 2>/dev/null || echo "WARN: oxlint install failed"; }
fi

echo "=== Recent commits ==="
git log --oneline -10 2>/dev/null || echo "(no commits yet)"

echo "=== Session started at $(date -u +"%Y-%m-%dT%H:%M:%SZ") ==="
