#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

: "${GO:=go}"
: "${PREFIX:=/tmp/ralleh-mcp-smoke-install}"
SHOP_PORT="${SHOP_PORT:-8621}"
TRAVEL_PORT="${TRAVEL_PORT:-8622}"
SHOP_ADDR="127.0.0.1:${SHOP_PORT}"
TRAVEL_ADDR="127.0.0.1:${TRAVEL_PORT}"

log() { printf '\n== %s ==\n' "$*"; }
cleanup() {
  if [[ -n "${SHOP_PID:-}" ]]; then kill "$SHOP_PID" 2>/dev/null || true; fi
  if [[ -n "${TRAVEL_PID:-}" ]]; then kill "$TRAVEL_PID" 2>/dev/null || true; fi
}
trap cleanup EXIT

log "go tests"
"$GO" test ./...

log "build"
GO="$GO" scripts/build.sh

log "one-shot health"
./dist/bin/ralleh-mcp-shop --health > /tmp/ralleh-mcp-smoke-shop-health.json
./dist/bin/ralleh-mcp-travel --health > /tmp/ralleh-mcp-smoke-travel-health.json
python3 - <<'PY'
import json
for name,path in [('shop','/tmp/ralleh-mcp-smoke-shop-health.json'),('travel','/tmp/ralleh-mcp-smoke-travel-health.json')]:
    data=json.load(open(path))
    assert data['ready'] is True, data
    assert data['status'] == 'ok', data
    assert data['collections'], data
    print(f"{name}: {data['service']} ready with {len(data['collections'])} collections")
PY

log "local HTTP health"
./dist/bin/ralleh-mcp-shop --health-server --health-listen "$SHOP_ADDR" >/tmp/ralleh-mcp-smoke-shop-server.log 2>&1 & SHOP_PID=$!
./dist/bin/ralleh-mcp-travel --health-server --health-listen "$TRAVEL_ADDR" >/tmp/ralleh-mcp-smoke-travel-server.log 2>&1 & TRAVEL_PID=$!
sleep 0.3
curl -fsS "http://${SHOP_ADDR}/healthz" >/tmp/ralleh-mcp-smoke-shop-healthz.json
curl -fsS "http://${SHOP_ADDR}/readyz" >/tmp/ralleh-mcp-smoke-shop-readyz.json
curl -fsS "http://${SHOP_ADDR}/version" >/tmp/ralleh-mcp-smoke-shop-version.json
curl -fsS "http://${TRAVEL_ADDR}/healthz" >/tmp/ralleh-mcp-smoke-travel-healthz.json
curl -fsS "http://${TRAVEL_ADDR}/readyz" >/tmp/ralleh-mcp-smoke-travel-readyz.json
curl -fsS "http://${TRAVEL_ADDR}/version" >/tmp/ralleh-mcp-smoke-travel-version.json
python3 - <<'PY'
import json
for name,path in [
    ('shop-healthz','/tmp/ralleh-mcp-smoke-shop-healthz.json'),
    ('shop-readyz','/tmp/ralleh-mcp-smoke-shop-readyz.json'),
    ('travel-healthz','/tmp/ralleh-mcp-smoke-travel-healthz.json'),
    ('travel-readyz','/tmp/ralleh-mcp-smoke-travel-readyz.json'),
]:
    data=json.load(open(path))
    assert data['ready'] is True, (name,data)
    print(f"{name}: {data['status']}")
for name,path in [
    ('shop-version','/tmp/ralleh-mcp-smoke-shop-version.json'),
    ('travel-version','/tmp/ralleh-mcp-smoke-travel-version.json'),
]:
    data=json.load(open(path))
    assert data.get('commit'), (name,data)
    assert data.get('version'), (name,data)
    print(f"{name}: {data}")
PY

log "local-only bind guard"
if ./dist/bin/ralleh-mcp-shop --health-server --health-listen 0.0.0.0:8623 >/tmp/ralleh-mcp-smoke-nonloop.out 2>/tmp/ralleh-mcp-smoke-nonloop.err; then
  cat /tmp/ralleh-mcp-smoke-nonloop.out
  echo "expected non-loopback bind to fail" >&2
  exit 1
fi
grep -q 'refusing non-loopback' /tmp/ralleh-mcp-smoke-nonloop.err
cat /tmp/ralleh-mcp-smoke-nonloop.err

log "install + backup flow"
rm -rf "$PREFIX"
PREFIX="$PREFIX" SRC=dist scripts/install.sh
PREFIX="$PREFIX" scripts/healthcheck.sh
BACKUP_PATH="$(PREFIX="$PREFIX" scripts/backup.sh | tail -1)"
test -f "$BACKUP_PATH"
tar -tzf "$BACKUP_PATH" >/tmp/ralleh-mcp-smoke-backup-list.txt
grep -q 'bin/ralleh-mcp-shop' /tmp/ralleh-mcp-smoke-backup-list.txt
grep -q 'configs/sources.shop.yaml' /tmp/ralleh-mcp-smoke-backup-list.txt
echo "backup ok: $BACKUP_PATH"

log "search capability boundary"
cat <<'MSG'
Actual product/flight source searches are intentionally not part of this smoke yet.
Current implementation has health, registries, budgets, affiliate safety, bounded execution, and ops hardening.
Real smoke searches should be added after MCP search handlers and at least one fake upstream adapter exist.
MSG

log "smoke PASS"
