#!/usr/bin/env bash
set -euo pipefail

input="$(cat)"
file="$(jq -r '.tool_input.file_path // .tool_input.path // empty' <<< "$input")"

# Go files
case "$file" in
  *.go)
    cd "$(git rev-parse --show-toplevel 2>/dev/null || dirname "$file")"
    if command -v golangci-lint &>/dev/null; then
      golangci-lint run --fix "$file" >/dev/null 2>&1 || true
      diag="$(golangci-lint run "$file" 2>&1 | head -20)"
      if [ -n "$diag" ]; then
        jq -Rn --arg msg "$diag" \
          '{ hookSpecificOutput: { hookEventName: "PostToolUse", additionalContext: $msg } }'
      fi
    fi
    ;;
  *.ts|*.tsx|*.js|*.jsx)
    npx biome format --write "$file" >/dev/null 2>&1 || true
    npx oxlint --fix "$file" >/dev/null 2>&1 || true
    diag="$(npx oxlint "$file" 2>&1 | head -20)"
    if [ -n "$diag" ]; then
      jq -Rn --arg msg "$diag" \
        '{ hookSpecificOutput: { hookEventName: "PostToolUse", additionalContext: $msg } }'
    fi
    ;;
esac
