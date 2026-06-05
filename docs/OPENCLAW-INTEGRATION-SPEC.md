# Ralleh MCP → OpenClaw Integration Spec

## Purpose

Define how each VPS-local `ralleh-mcp` install should be consumed by OpenClaw running on that same VPS.

This document is intentionally precise about **current state vs target state**.

## Reality check: current state

Today, this repo ships local service binaries with health endpoints and deterministic smoke behavior:

- `ralleh-mcp-shop`
- `ralleh-mcp-travel`
- `ralleh-mcp-search`
- `ralleh-mcp-brand`

Current MCP readiness:

- `ralleh-mcp-search --mcp` is the first stdio MCP entrypoint and exposes a small content-search tool surface.
- `shop`, `travel`, and `brand` are still local service binaries, not MCP protocol servers yet.

That means:

- this repo is already valid as a **VPS-local toolbelt**;
- `ralleh-mcp-search` is the first binary ready for OpenClaw `mcp.servers` validation;
- remaining binaries should not be registered as MCP servers until protocol handlers are added.

## Target model

For each VPS:

1. `ralleh-mcp` is installed locally under `/opt/ralleh/ralleh-mcp`.
2. OpenClaw runs on the same VPS.
3. OpenClaw consumes Ralleh MCP capabilities through **local-only MCP server definitions**.
4. No Docker is required.
5. No public exposure is required by default.

## Preferred OpenClaw consumption model

### End state

OpenClaw should consume Ralleh MCP through **local stdio MCP servers** registered in `mcp.servers`.

Why stdio first:

- simplest local trust boundary;
- no extra network listener required for tool calls;
- aligns with OpenClaw's native MCP registry model for local child processes;
- keeps the VPS-local architecture tight and private.

### Preferred registration shape

Target config shape in OpenClaw after each binary has an MCP stdio mode. Today only `ralleh-search` should be treated as implementation-ready:

```json
{
  "mcp": {
    "servers": {
      "ralleh-shop": {
        "command": "/opt/ralleh/ralleh-mcp/bin/ralleh-mcp-shop",
        "args": []
      },
      "ralleh-travel": {
        "command": "/opt/ralleh/ralleh-mcp/bin/ralleh-mcp-travel",
        "args": []
      },
      "ralleh-search": {
        "command": "/opt/ralleh/ralleh-mcp/bin/ralleh-mcp-search",
        "args": ["--mcp"]
      },
      "ralleh-brand": {
        "command": "/opt/ralleh/ralleh-mcp/bin/ralleh-mcp-brand",
        "args": ["--db", "/opt/ralleh/ralleh-mcp/brand.db"]
      }
    }
  }
}
```

Important: `ralleh-search` is the current working pattern. The other definitions are target-state until those binaries speak MCP correctly over stdio.

## OpenClaw tool-policy requirements

When OpenClaw consumes saved MCP servers, their tools are exposed under the `bundle-mcp` plugin surface.

So the VPS-local OpenClaw runtime must allow MCP/plugin tools through normal tool policy.

### Minimum policy considerations

- tool profile must not be too restrictive;
- if sandbox tool gating is active, allow `bundle-mcp` or `group:plugins` at the sandbox layer;
- if using restrictive plugin policy, ensure the MCP/plugin surface is allowed.

### Practical recommendation

For agent/tool surfaces that should use local Ralleh MCP tools:

- use at least a tool profile compatible with plugin/MCP tool usage;
- ensure sandbox tool allowlist includes `bundle-mcp` when sandbox mode filters plugin tools;
- do not expose these tools to agents/senders that do not need them.

Example sandbox-layer allow shape:

```json
{
  "tools": {
    "sandbox": {
      "tools": {
        "alsoAllow": ["bundle-mcp"]
      }
    }
  }
}
```

## Local-only security model

The integration must preserve the VPS-local design:

- `ralleh-mcp` binaries run on the same host as OpenClaw;
- health endpoints remain loopback-only;
- MCP execution should be local child-process stdio, not public network services;
- no Docker runtime dependency;
- no implicit widening to `0.0.0.0` just to make tooling easier.

## Brand DB contract

`ralleh-mcp-brand` requires a local DB path.

Default contract:

```text
/opt/ralleh/ralleh-mcp/brand.db
```

If a VPS chooses another local path, the value must stay consistent across:

- systemd units;
- install/health scripts;
- OpenClaw MCP server args.

## Current `ralleh-mcp-search --mcp` tool surface

`ralleh-mcp-search --mcp` exposes:

- `list_collections`
- `rank_sources`
- `search_content`

The first implementation intentionally uses deterministic fake adapters, matching the repo's smoke-test posture. Live website/API extraction remains a later adapter milestone.

## What is missing before all binaries are direct OpenClaw MCP registrations

Each remaining intended MCP binary still needs:

1. **real MCP protocol handlers** over stdio;
2. deterministic tool registration and schema definitions;
3. stable tool names per server;
4. request/response handling that matches MCP expectations;
5. integration tests proving OpenClaw can launch and use the local MCP server.

## Recommended implementation sequence

### Phase 1 — Search MCP baseline

`ralleh-mcp-search --mcp` now provides the baseline stdio MCP implementation.

Remaining validation for this phase:

- register it in a real OpenClaw `mcp.servers` config;
- confirm tool visibility through OpenClaw tool policy;
- run one OpenClaw-driven tool call.

### Phase 2 — Add `shop`

After `search` works end-to-end:

- add MCP protocol support for `ralleh-mcp-shop`;
- keep strict research-only boundaries;
- keep affiliate handling transparent and non-ranking by default.

### Phase 3 — Add `travel`

Only after shop/search are stable:

- add MCP protocol support for `ralleh-mcp-travel`;
- keep no-booking / no-card / no-PII boundaries hard.

### Phase 4 — Add `brand`

`ralleh-mcp-brand` is useful, but persistence and write paths make it more sensitive.

Before broad agent exposure:

- finalize write/read audit contract;
- ensure DB path and permissions are explicit;
- scope tool exposure carefully.

## Recommended first OpenClaw-facing tool shapes

These are conceptual, not implemented names.

### `ralleh-mcp-search`

Good first tool candidates:

- `search_content`
- `rank_sources`
- `list_collections`

### `ralleh-mcp-shop`

Good first tool candidates:

- `search_products`
- `rank_sources`
- `list_collections`

### `ralleh-mcp-travel`

Good first tool candidates:

- `search_flights`
- `rank_sources`
- `list_collections`

### `ralleh-mcp-brand`

Good first tool candidates:

- `get_brand_profile`
- `get_brand_voice`
- `validate_content`

Do not expose write/mutation tools broadly until role and audit boundaries are settled.

## Verification contract for the first real integration

A server is not considered OpenClaw-ready until all of these are true:

1. local binary still passes repo tests;
2. MCP mode starts correctly as a stdio child process;
3. OpenClaw can register it under `mcp.servers`;
4. the resulting MCP tools are visible in an allowed tool policy;
5. at least one real OpenClaw-driven tool call succeeds;
6. no Docker dependency was introduced;
7. no public listener was required.

## Non-goals

- no Docker Toolkit;
- no remote shared MCP gateway for all VPSes in v1;
- no public internet MCP exposure by default;
- no pretending health endpoints are MCP transport endpoints.

## Bottom line

The correct architecture is:

- **Ralleh MCP installed locally per VPS**
- **OpenClaw on the same VPS**
- **local stdio MCP integration once protocol handlers exist**

The correct immediate next engineering move is:

- prove OpenClaw can consume `ralleh-mcp-search --mcp` locally;
- then expand server-by-server.
