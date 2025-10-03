#!/bin/bash
set -e

COMMAND=${1:-up}
MIGRATE_BIN="./bin/migrate"

echo "��� Corporate API - Database Migration"
echo "====================================="

# Check .env
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found"
    exit 1
fi

# Load environment (.env must be simple KEY=VALUE lines)
if command -v bash >/dev/null 2>&1; then
    # safer for bash: export all variables sourced from .env
    set -o allexport
    # shellcheck disable=SC1091
    source .env
    set +o allexport
else
    # fallback (less safe for values with spaces)
    export $(cat .env | grep -v '^#' | xargs)
fi

# Build if not exists
if [ ! -f "$MIGRATE_BIN" ]; then
    echo "Building migration tool..."
    mkdir -p bin
    go build -o "$MIGRATE_BIN" cmd/migrate/main.go
    echo "✅ Migration tool built"
fi

# Run command
case "$COMMAND" in
    up)
        echo "��� Running migrations..."
        $MIGRATE_BIN -up
        ;;
    down)
        echo "⏮️  Rolling back..."
        $MIGRATE_BIN -down
        ;;
    status)
        echo "��� Checking status..."
        $MIGRATE_BIN -status
        ;;
    *)
        echo "Usage: $0 [up|down|status]"
        exit 1
        ;;
esac

echo "✅ Done!"
