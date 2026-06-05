#!/usr/bin/env bash
set -euo pipefail
PREFIX="${PREFIX:-/opt/ralleh/ralleh-mcp}"
BRAND_DB="${BRAND_DB:-$PREFIX/brand.db}"
"$PREFIX/bin/ralleh-mcp-shop" --health >/dev/null
"$PREFIX/bin/ralleh-mcp-travel" --health >/dev/null
"$PREFIX/bin/ralleh-mcp-search" --health >/dev/null
"$PREFIX/bin/ralleh-mcp-brand" --db "$BRAND_DB" --health >/dev/null
echo "Ralleh MCP health OK"
