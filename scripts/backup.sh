#!/usr/bin/env bash
set -euo pipefail
PREFIX="${PREFIX:-/opt/ralleh/ralleh-mcp}"
BACKUP_DIR="${BACKUP_DIR:-$PREFIX/backups}"
BRAND_DB="${BRAND_DB:-$PREFIX/brand.db}"
STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
ARCHIVE="$BACKUP_DIR/ralleh-mcp-$STAMP.tgz"

mkdir -p "$BACKUP_DIR"

paths=(bin configs)
if [[ -f "$BRAND_DB" ]]; then
  paths+=("$(basename "$BRAND_DB")")
fi

tar -C "$PREFIX" -czf "$ARCHIVE" "${paths[@]}"
echo "$ARCHIVE"
