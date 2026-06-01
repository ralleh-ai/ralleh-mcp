#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

: "${GO:=go}"
: "${PREFIX:=/tmp/ralleh-mcp-smoke-install}"
SHOP_PORT="${SHOP_PORT:-8621}"
TRAVEL_PORT="${TRAVEL_PORT:-8622}"
SEARCH_PORT="${SEARCH_PORT:-8624}"
SHOP_ADDR="127.0.0.1:${SHOP_PORT}"
TRAVEL_ADDR="127.0.0.1:${TRAVEL_PORT}"
SEARCH_ADDR="127.0.0.1:${SEARCH_PORT}"

log() { printf '\n== %s ==\n' "$*"; }
cleanup() {
  if [[ -n "${SHOP_PID:-}" ]]; then kill "$SHOP_PID" 2>/dev/null || true; fi
  if [[ -n "${TRAVEL_PID:-}" ]]; then kill "$TRAVEL_PID" 2>/dev/null || true; fi
  if [[ -n "${SEARCH_PID:-}" ]]; then kill "$SEARCH_PID" 2>/dev/null || true; fi
}
trap cleanup EXIT

log "go tests"
"$GO" test ./...

log "build"
GO="$GO" scripts/build.sh

log "one-shot health"
./dist/bin/ralleh-mcp-shop --health > /tmp/ralleh-mcp-smoke-shop-health.json
./dist/bin/ralleh-mcp-travel --health > /tmp/ralleh-mcp-smoke-travel-health.json
./dist/bin/ralleh-mcp-search --health > /tmp/ralleh-mcp-smoke-search-health.json
python3 - <<'PY'
import json
for name,path in [('shop','/tmp/ralleh-mcp-smoke-shop-health.json'),('travel','/tmp/ralleh-mcp-smoke-travel-health.json'),('search','/tmp/ralleh-mcp-smoke-search-health.json')]:
    data=json.load(open(path))
    assert data['ready'] is True, data
    assert data['status'] == 'ok', data
    assert data['collections'], data
    print(f"{name}: {data['service']} ready with {len(data['collections'])} collections")
PY

log "local HTTP health"
./dist/bin/ralleh-mcp-shop --health-server --health-listen "$SHOP_ADDR" >/tmp/ralleh-mcp-smoke-shop-server.log 2>&1 & SHOP_PID=$!
./dist/bin/ralleh-mcp-travel --health-server --health-listen "$TRAVEL_ADDR" >/tmp/ralleh-mcp-smoke-travel-server.log 2>&1 & TRAVEL_PID=$!
./dist/bin/ralleh-mcp-search --health-server --health-listen "$SEARCH_ADDR" >/tmp/ralleh-mcp-smoke-search-server.log 2>&1 & SEARCH_PID=$!
sleep 0.3
curl -fsS "http://${SHOP_ADDR}/healthz" >/tmp/ralleh-mcp-smoke-shop-healthz.json
curl -fsS "http://${SHOP_ADDR}/readyz" >/tmp/ralleh-mcp-smoke-shop-readyz.json
curl -fsS "http://${SHOP_ADDR}/version" >/tmp/ralleh-mcp-smoke-shop-version.json
curl -fsS "http://${TRAVEL_ADDR}/healthz" >/tmp/ralleh-mcp-smoke-travel-healthz.json
curl -fsS "http://${TRAVEL_ADDR}/readyz" >/tmp/ralleh-mcp-smoke-travel-readyz.json
curl -fsS "http://${TRAVEL_ADDR}/version" >/tmp/ralleh-mcp-smoke-travel-version.json
curl -fsS "http://${SEARCH_ADDR}/healthz" >/tmp/ralleh-mcp-smoke-search-healthz.json
curl -fsS "http://${SEARCH_ADDR}/readyz" >/tmp/ralleh-mcp-smoke-search-readyz.json
curl -fsS "http://${SEARCH_ADDR}/version" >/tmp/ralleh-mcp-smoke-search-version.json
python3 - <<'PY'
import json
for name,path in [
    ('shop-healthz','/tmp/ralleh-mcp-smoke-shop-healthz.json'),
    ('shop-readyz','/tmp/ralleh-mcp-smoke-shop-readyz.json'),
    ('travel-healthz','/tmp/ralleh-mcp-smoke-travel-healthz.json'),
    ('travel-readyz','/tmp/ralleh-mcp-smoke-travel-readyz.json'),
    ('search-healthz','/tmp/ralleh-mcp-smoke-search-healthz.json'),
    ('search-readyz','/tmp/ralleh-mcp-smoke-search-readyz.json'),
]:
    data=json.load(open(path))
    assert data['ready'] is True, (name,data)
    print(f"{name}: {data['status']}")
for name,path in [
    ('shop-version','/tmp/ralleh-mcp-smoke-shop-version.json'),
    ('travel-version','/tmp/ralleh-mcp-smoke-travel-version.json'),
    ('search-version','/tmp/ralleh-mcp-smoke-search-version.json'),
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

log "deterministic fake search smoke"
./dist/bin/ralleh-mcp-shop --search-query "cordless drill" --search-collection tools --search-sources ebay,random_site >/tmp/ralleh-mcp-smoke-shop-search.json
./dist/bin/ralleh-mcp-travel --flight-origin MCO --flight-destination LAS --flight-depart 2026-07-12 --flight-collection us_domestic_flights --flight-sources duffel,random_ota >/tmp/ralleh-mcp-smoke-travel-search.json
./dist/bin/ralleh-mcp-search --search-query "ai chips" --search-collection technology --search-sources hacker_news,random_blog >/tmp/ralleh-mcp-smoke-content-search.json
python3 - <<'PY'
import json
shop=json.load(open('/tmp/ralleh-mcp-smoke-shop-search.json'))
travel=json.load(open('/tmp/ralleh-mcp-smoke-travel-search.json'))
content=json.load(open('/tmp/ralleh-mcp-smoke-content-search.json'))
assert shop['status']=='ok', shop
assert travel['status']=='ok', travel
assert content['status']=='ok', content
assert shop['results'], shop
assert travel['results'], travel
assert content['results'], content
assert 'random_site' in shop['sourcePlan']['rejectedSources'], shop['sourcePlan']
assert 'random_ota' in travel['sourcePlan']['rejectedSources'], travel['sourcePlan']
assert 'random_blog' in content['sourcePlan']['rejectedSources'], content['sourcePlan']
assert any(item['affiliate']['applied'] for item in shop['results']), shop['results']
assert travel['capabilities']['canBook'] is False, travel['capabilities']
assert travel['capabilities']['canUseCreditCard'] is False, travel['capabilities']
assert content['capabilities']['canCrawlArbitraryWebsites'] is False, content['capabilities']
print(f"shop fake search: {len(shop['results'])} results, rejected={shop['sourcePlan']['rejectedSources']}")
print(f"travel fake search: {len(travel['results'])} results, rejected={travel['sourcePlan']['rejectedSources']}")
print(f"content fake search: {len(content['results'])} results, rejected={content['sourcePlan']['rejectedSources']}")
PY
cat <<'MSG'
Fake upstream searches validate search contracts, source rejection, affiliate URL path, no-booking/no-card boundaries, and no-arbitrary-crawl content search.
Live website/API searches are still intentionally not part of smoke until real source adapters are implemented.
MSG

log "smoke PASS"
