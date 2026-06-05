# Ralleh MCP VPS Rollout Checklist

Use this when deploying or upgrading `ralleh-mcp` on a VPS that also runs OpenClaw.

## 1) Preflight

- Confirm repo checkout path is correct.
- Confirm Go toolchain is available.
- Confirm target install prefix (default: `/opt/ralleh/ralleh-mcp`).
- Confirm service account/user exists if using systemd units (`ralleh-mcp`).
- Confirm health ports are free on localhost:
  - `8621` shop
  - `8622` travel
  - `8624` search
  - `8625` brand

## 2) Build and verify in repo

```bash
go test ./...
scripts/build.sh
scripts/smoke.sh
```

Expected result: tests and smoke both pass.

## 3) Backup current install

```bash
PREFIX=/opt/ralleh/ralleh-mcp scripts/backup.sh
```

Expected result: timestamped tarball under `/opt/ralleh/ralleh-mcp/backups/`.

## 4) Install upgrade

```bash
PREFIX=/opt/ralleh/ralleh-mcp SRC=dist scripts/install.sh
PREFIX=/opt/ralleh/ralleh-mcp scripts/healthcheck.sh
```

If brand state should live somewhere else, override `BRAND_DB` consistently during install/health checks and in the brand systemd unit.

## 5) Install or refresh systemd units

Copy these units from `deploy/systemd/` into `/etc/systemd/system/`:

- `ralleh-mcp-shop-health.service`
- `ralleh-mcp-travel-health.service`
- `ralleh-mcp-search-health.service`
- `ralleh-mcp-brand-health.service`

Then reload and enable:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now ralleh-mcp-shop-health.service
sudo systemctl enable --now ralleh-mcp-travel-health.service
sudo systemctl enable --now ralleh-mcp-search-health.service
sudo systemctl enable --now ralleh-mcp-brand-health.service
```

## 6) Local verification

```bash
curl -fsS http://127.0.0.1:8621/healthz
curl -fsS http://127.0.0.1:8622/readyz
curl -fsS http://127.0.0.1:8624/version
curl -fsS http://127.0.0.1:8625/healthz
```

Expected result: JSON responses with `ready: true` where applicable.

## 7) OpenClaw integration verification

- Confirm OpenClaw on the same VPS is configured to use the local MCP services/binaries.
- Confirm no Docker dependency is referenced in install/runtime docs.
- Confirm services remain bound to `127.0.0.1` unless there is an explicit operator decision otherwise.
- Confirm no transaction/book/checkout/card behavior is exposed.

## 8) Post-upgrade checks

- `systemctl status` is clean for enabled health services.
- Backup tarball exists and is restorable.
- Brand DB path is correct and writable for the brand unit.
- Repo working tree is clean after deployment-specific edits are handled.

## 9) If something fails

- Stop and inspect the relevant health unit logs.
- Re-run one-shot health checks from the installed binaries.
- Restore from the latest backup tarball if needed.
- Do not widen network exposure just to make health checks pass.
