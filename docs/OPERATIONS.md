# Ralleh MCP Operations (VPS-local OpenClaw model)

## Operating model

Ralleh MCP is run **locally on each VPS** where OpenClaw runs.

- install binaries to `/opt/ralleh/ralleh-mcp`;
- keep health endpoints on loopback by default;
- do not require Docker for install or runtime.

## VPS-local security

- Health endpoints bind to `127.0.0.1` by default.
- Non-loopback health binding is rejected unless `--allow-non-loopback-health` is explicitly set.
- Systemd units include hardening options and no ambient privileges.
- No v1 tools perform booking, checkout, credit card handling, account login, or passenger PII collection.

## Health checks

One-shot:

```bash
ralleh-mcp-shop --health
ralleh-mcp-travel --health
ralleh-mcp-search --health
ralleh-mcp-brand --db /opt/ralleh/ralleh-mcp/brand.db --health
```

Local HTTP:

```bash
ralleh-mcp-shop --health-server --health-listen 127.0.0.1:8621
ralleh-mcp-travel --health-server --health-listen 127.0.0.1:8622
ralleh-mcp-search --health-server --health-listen 127.0.0.1:8624
ralleh-mcp-brand --db /opt/ralleh/ralleh-mcp/brand.db --health-server --health-listen 127.0.0.1:8625

curl -fsS http://127.0.0.1:8621/healthz
curl -fsS http://127.0.0.1:8622/readyz
curl -fsS http://127.0.0.1:8624/version
curl -fsS http://127.0.0.1:8625/healthz
```

## Install layout

```text
/opt/ralleh/ralleh-mcp/
  bin/
    ralleh-mcp-shop
    ralleh-mcp-travel
    ralleh-mcp-search
    ralleh-mcp-brand
  configs/
    sources.shop.yaml
    sources.travel.yaml
  backups/
```

Default brand DB path used by ops scripts/units:

```text
/opt/ralleh/ralleh-mcp/brand.db
```

(Override with `--db` if you want a different local path such as `/var/lib/...`.)

## Systemd units

Unit templates are in `deploy/systemd/` for local health services:

- `ralleh-mcp-shop-health.service`
- `ralleh-mcp-travel-health.service`
- `ralleh-mcp-search-health.service`
- `ralleh-mcp-brand-health.service`

After copying units into `/etc/systemd/system/`:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now ralleh-mcp-shop-health.service
sudo systemctl enable --now ralleh-mcp-travel-health.service
sudo systemctl enable --now ralleh-mcp-search-health.service
sudo systemctl enable --now ralleh-mcp-brand-health.service
```

## Upgrade model

1. Build new binaries.
2. Run `scripts/backup.sh`.
3. Install binaries/configs using `scripts/install.sh`.
4. Run `scripts/healthcheck.sh`.
5. Restart affected systemd units.
6. Verify `/healthz` and `/readyz` locally.

## Backup model

`scripts/backup.sh` archives `bin/`, `configs/`, and (if present) `brand.db` from `$PREFIX`.

This keeps local VPS state recoverable without adding external dependencies.
