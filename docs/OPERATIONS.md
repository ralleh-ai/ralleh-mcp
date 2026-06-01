# Ralleh MCP Operations

## VPS-local security

Ralleh MCP must run private to the VPS/OpenClaw runtime by default.

- Health endpoints bind to `127.0.0.1` by default.
- Non-loopback health binding is rejected unless `--allow-non-loopback-health` is explicitly set.
- Systemd units include hardening options and no ambient privileges.
- No v1 tools perform booking, checkout, credit card handling, account login, or passenger PII collection.

## Health checks

One-shot:

```bash
ralleh-mcp-shop --health
ralleh-mcp-travel --health
```

Local HTTP:

```bash
ralleh-mcp-shop --health-server --health-listen 127.0.0.1:8621
curl -fsS http://127.0.0.1:8621/healthz
curl -fsS http://127.0.0.1:8621/readyz
curl -fsS http://127.0.0.1:8621/version
```

## Install layout

```text
/opt/ralleh/ralleh-mcp/
  bin/
    ralleh-mcp-shop
    ralleh-mcp-travel
  configs/
    sources.shop.yaml
    sources.travel.yaml
  backups/
```

## Upgrade model

1. Build new binaries.
2. Run `scripts/backup.sh`.
3. Install binaries/configs using `scripts/install.sh`.
4. Run `scripts/healthcheck.sh`.
5. Restart affected systemd units.
6. Verify `/healthz` and `/readyz` locally.

## Backup model

Back up configs and existing binaries before upgrade. Search cache/state is intentionally not required yet; when circuit-breaker/cache persistence lands, add it to the backup manifest.
