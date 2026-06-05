# Ralleh MCP VPS-Local Toolbelt Spec

## Purpose

This repo is the **VPS-local MCP toolbelt** for Ralleh.

Deployment model:

- each VPS has its own local checkout/build/install of this repo;
- OpenClaw instances running on that same VPS invoke these MCP binaries locally;
- services stay private to localhost unless an operator explicitly opts out.

## Scope (current)

Shipped binaries:

- `ralleh-mcp-shop` (shopping research)
- `ralleh-mcp-travel` (travel research)
- `ralleh-mcp-search` (content research)
- `ralleh-mcp-brand` (brand memory + validation)

Shared properties:

- curated source/collection model (no arbitrary crawling);
- bounded runtime budgets;
- local health checks (`--health` and optional localhost HTTP health server);
- no transaction execution (no checkout/booking/credit card handling).

## Explicit non-goals

- No Docker packaging or Docker runtime requirement in this repo.
- No remote multi-tenant control plane in this repo.
- No public internet health exposure by default.

## Operational contract

Install target (default):

```text
/opt/ralleh/ralleh-mcp/
  bin/
  configs/
  backups/
```

Default local brand DB path in shipped ops artifacts:

```text
/opt/ralleh/ralleh-mcp/brand.db
```

Operators can override to another local path (for example `/var/lib/...`).

Systemd artifacts are local-host oriented and should bind health endpoints to `127.0.0.1`.

## Reliability and safety baseline

- Keep docs aligned with real binaries and scripts.
- Prefer small reversible ops changes over speculative architecture.
- Verify changes with `go test ./...` and `scripts/smoke.sh` when feasible.

## Near-term roadmap (repo-local)

1. Keep README/OPERATIONS/ARCHITECTURE anchored to VPS-local OpenClaw usage.
2. Maintain systemd unit parity for every shipped service binary.
3. Keep backup scripts aligned with persisted local state (including brand DB when configured).
4. Add only incremental hardening that is testable in-repo.
