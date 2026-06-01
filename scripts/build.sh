#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
: "${GO:=go}"
VERSION="${VERSION:-dev}"
COMMIT="${COMMIT:-$(git rev-parse --short=12 HEAD 2>/dev/null || echo unknown)}"
DATE="${DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"
LDFLAGS="-s -w -X github.com/ralleh-ai/ralleh-mcp/internal/core/runtime.Version=$VERSION -X github.com/ralleh-ai/ralleh-mcp/internal/core/runtime.Commit=$COMMIT -X github.com/ralleh-ai/ralleh-mcp/internal/core/runtime.Date=$DATE"
mkdir -p dist/bin
"$GO" test ./...
"$GO" build -trimpath -ldflags "$LDFLAGS" -o dist/bin/ralleh-mcp-shop ./cmd/ralleh-mcp-shop
"$GO" build -trimpath -ldflags "$LDFLAGS" -o dist/bin/ralleh-mcp-travel ./cmd/ralleh-mcp-travel
"$GO" build -trimpath -ldflags "$LDFLAGS" -o dist/bin/ralleh-mcp-search ./cmd/ralleh-mcp-search
"$GO" build -trimpath -ldflags "$LDFLAGS" -o dist/bin/ralleh-mcp-brand ./cmd/ralleh-mcp-brand
sha256sum dist/bin/ralleh-mcp-shop dist/bin/ralleh-mcp-travel dist/bin/ralleh-mcp-search dist/bin/ralleh-mcp-brand > dist/SHA256SUMS
