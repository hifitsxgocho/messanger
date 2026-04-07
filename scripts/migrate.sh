#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/../backend"

DIRECTION="${1:-up}"
STEPS="${2:-0}"

export DATABASE_URL="${DATABASE_URL:-postgres://messenger:messenger_secret@localhost:5432/messenger?sslmode=disable}"
export MIGRATIONS_PATH="${MIGRATIONS_PATH:-migrations}"

echo "Database : $DATABASE_URL"
echo "Direction: $DIRECTION"
echo ""

cd "$BACKEND_DIR"
go run ./cmd/migrate -direction="$DIRECTION" -steps="$STEPS"
