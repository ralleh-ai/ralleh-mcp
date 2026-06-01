#!/usr/bin/env bash
set -euo pipefail
PREFIX="${PREFIX:-/opt/ralleh/ralleh-mcp}"
"$PREFIX/bin/ralleh-mcp-shop" --health >/dev/null
"$PREFIX/bin/ralleh-mcp-travel" --health >/dev/null
"$PREFIX/bin/ralleh-mcp-search" --health >/dev/null
echo "Ralleh MCP health OK"
