#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
: "${REGISTRY_DEPLOY_PASSWORD:?Set REGISTRY_DEPLOY_PASSWORD before exporting the gated GitHub Pages build}"
exec go run ./cmd/registry export --out docs "$@"
