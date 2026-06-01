# Ralleh MCP

Enterprise-grade MCP services for bounded, curated shopping and travel research.

## Mission

Ralleh MCP gives OpenClaw and other LLM clients fast, structured, safe search tools for:

- shopping/product research via `ralleh-mcp-shop`
- travel/flight research via `ralleh-mcp-travel`

The services are **research engines**, not transaction agents.

## Non-negotiable v1 boundaries

- No credit cards.
- No checkout.
- No booking.
- No passenger PII.
- No account login automation.
- No arbitrary website crawling.
- No captcha bypass.
- No unbounded requests, goroutines, retries, browser sessions, or response bodies.

## Core architecture

```text
cmd/
  ralleh-mcp-shop/       # shopping MCP server
  ralleh-mcp-travel/     # travel MCP server
internal/
  core/                  # shared enterprise request/runtime primitives
  shop/                  # curated shopping collections, adapters, ranking
  travel/                # curated travel/flight collections, adapters, ranking
configs/
  sources.shop.yaml      # curated shopping source registry
  sources.travel.yaml    # curated travel source registry
```

## Design

The LLM chooses a **known collection**, not random URLs:

```json
{
  "query": "cordless drill brushless 20v",
  "collection": "tools",
  "preferredSources": ["home_depot", "lowes", "harbor_freight"],
  "budgetProfile": "standard"
}
```

Ralleh MCP then:

1. validates the collection and known source IDs;
2. resolves approved adapters from curated registries;
3. clamps the request to hard resource budgets;
4. searches in bounded parallel workers;
5. returns normalized results, source diagnostics, evidence, and affiliate-ready presentation URLs.

## Affiliate URL rule

Ralleh MCP keeps two URLs:

- `canonicalUrl` for fetch, cache, dedupe, and evidence;
- `presentedUrl` for the user-facing affiliate-tagged link.

Affiliate status must never silently affect ranking unless a ranking policy explicitly allows it.

## Current status

Initial scaffold. The first production target is `shop.search` with curated source collections and no browser fallback.
