# Ralleh MCP Architecture

## Deployment shape

Ralleh MCP is designed as a **VPS-local toolbelt**:

- each VPS has local binaries and configs;
- OpenClaw on that VPS calls local MCP services;
- health/ops endpoints are localhost-first.

No Docker dependency is required for this architecture.

## Services

- `ralleh-mcp-shop`: curated shopping/product research.
- `ralleh-mcp-travel`: curated travel/flight research.
- `ralleh-mcp-search`: curated content/news research.
- `ralleh-mcp-brand`: local brand memory and compliance validation.
- `internal/core`: shared budget, source, affiliate, result, request, cache, and observability primitives.

## Boundary

Ralleh MCP accepts collection/source IDs, not arbitrary website URLs. This prevents accidental crawler behavior and lets the service maintain known adapters, budgets, and health state.

## Request safety

Every search must have:

- global deadline
- per-source deadline
- max concurrency
- max source count
- max result count
- max response bytes
- retry budget
- source circuit breaker state before production launch

## Captchas and bot challenges

The goal is browser-compatible fetching, not bypass. Captchas/challenges are detected, reported, and treated as blocked/degraded source states. Human-in-the-loop browser verification can be added later.

## Transactions

No v1 transaction tools exist. The service returns research, recommendations, warnings, and user-facing links only.
