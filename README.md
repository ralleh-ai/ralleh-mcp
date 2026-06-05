# Ralleh MCP (VPS-local OpenClaw toolbelt)

Ralleh MCP is a **VPS-local MCP toolbelt repo**.

Each VPS runs its own local Ralleh MCP binaries, and OpenClaw instances on that same VPS call them locally.

> This repo explicitly does **not** adopt Docker as part of its runtime model.

See the target-role spec: [`docs/VPS-LOCAL-TOOLBELT-SPEC.md`](docs/VPS-LOCAL-TOOLBELT-SPEC.md).

## Services in this repo

- `ralleh-mcp-shop` — curated shopping research
- `ralleh-mcp-travel` — curated travel/flight research
- `ralleh-mcp-search` — curated content/news research
- `ralleh-mcp-brand` — local brand memory + validation

These are research/memory tools, not transaction agents.

## Hard boundaries

- No credit cards
- No checkout
- No booking
- No passenger PII handling
- No account login automation
- No arbitrary website crawling
- No captcha bypass

## Local-first operations

One-shot health checks:

```bash
ralleh-mcp-shop --health
ralleh-mcp-travel --health
ralleh-mcp-search --health
ralleh-mcp-brand --db /opt/ralleh/ralleh-mcp/brand.db --health
```

Local-only HTTP health endpoints:

```bash
ralleh-mcp-shop --health-server --health-listen 127.0.0.1:8621
ralleh-mcp-travel --health-server --health-listen 127.0.0.1:8622
ralleh-mcp-search --health-server --health-listen 127.0.0.1:8624
ralleh-mcp-brand --db /opt/ralleh/ralleh-mcp/brand.db --health-server --health-listen 127.0.0.1:8625
```

Non-loopback health bind is rejected unless `--allow-non-loopback-health` is explicitly set.

## Install layout

```text
/opt/ralleh/ralleh-mcp/
  bin/
  configs/
  backups/
```

Default brand DB path used by ops scripts/units:

```text
/opt/ralleh/ralleh-mcp/brand.db
```

(You can override `--db` if you prefer `/var/lib/...`.)

## Build, test, smoke

```bash
go test ./...
scripts/build.sh
scripts/smoke.sh
```

## Operations docs

- [`docs/OPERATIONS.md`](docs/OPERATIONS.md) — install/upgrade/backup/systemd
- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — service and safety model
- [`docs/VPS-LOCAL-TOOLBELT-SPEC.md`](docs/VPS-LOCAL-TOOLBELT-SPEC.md) — repo role and roadmap
