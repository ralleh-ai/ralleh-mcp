#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
scripts/build.sh
scripts/backup.sh
scripts/install.sh
scripts/healthcheck.sh
