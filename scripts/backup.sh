#!/usr/bin/env bash
set -euo pipefail
PREFIX="${PREFIX:-/opt/ralleh/ralleh-mcp}"
BACKUP_DIR="${BACKUP_DIR:-$PREFIX/backups}"
STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
mkdir -p "$BACKUP_DIR"
tar -C "$PREFIX" -czf "$BACKUP_DIR/ralleh-mcp-$STAMP.tgz" bin configs 2>/dev/null || true
echo "$BACKUP_DIR/ralleh-mcp-$STAMP.tgz"
