#!/usr/bin/env bash
set -euo pipefail
PREFIX="${PREFIX:-/opt/ralleh/ralleh-mcp}"
SRC="${SRC:-dist}"
install -d -m 0755 "$PREFIX/bin" "$PREFIX/configs" "$PREFIX/backups"
install -m 0755 "$SRC/bin/ralleh-mcp-shop" "$PREFIX/bin/ralleh-mcp-shop"
install -m 0755 "$SRC/bin/ralleh-mcp-travel" "$PREFIX/bin/ralleh-mcp-travel"
install -m 0755 "$SRC/bin/ralleh-mcp-search" "$PREFIX/bin/ralleh-mcp-search"
install -m 0755 "$SRC/bin/ralleh-mcp-brand" "$PREFIX/bin/ralleh-mcp-brand"
install -m 0644 configs/sources.shop.yaml "$PREFIX/configs/sources.shop.yaml"
install -m 0644 configs/sources.travel.yaml "$PREFIX/configs/sources.travel.yaml"
"$PREFIX/bin/ralleh-mcp-shop" --health >/dev/null
"$PREFIX/bin/ralleh-mcp-travel" --health >/dev/null
"$PREFIX/bin/ralleh-mcp-search" --health >/dev/null
"$PREFIX/bin/ralleh-mcp-brand" --db "$PREFIX/brand-smoke.db" --health >/dev/null
echo "Installed Ralleh MCP to $PREFIX"
