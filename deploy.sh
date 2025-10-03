#!/usr/bin/env bash
set -euo pipefail

COMPOSE_FILE=docker-compose.prod.yml

echo "== Deploy script: pull images, migrate, and start api"

if [ ! -f .env.prod ]; then
  echo "Error: .env.prod not found. Create it from .env and set production values." >&2
  exit 1
fi

# Pull fresh images
echo "Pulling images..."
docker compose -f "$COMPOSE_FILE" pull

# Start postgres and redis (but don't start api yet)
echo "Starting postgres and redis..."
docker compose -f "$COMPOSE_FILE" up -d postgres redis

echo "Waiting for Postgres to be ready..."
RETRIES=30
until docker compose -f "$COMPOSE_FILE" exec -T postgres pg_isready -U "${DB_USER}" >/dev/null 2>&1; do
  RETRIES=$((RETRIES-1))
  if [ $RETRIES -le 0 ]; then
    echo "Postgres did not become ready in time." >&2
    docker compose -f "$COMPOSE_FILE" logs --no-color postgres | tail -n 50
    exit 1
  fi
  sleep 2
done

echo "Running migrations (migrator)..."
# Run migrator as one-off; it should use the same image and exit 0 on success
docker compose -f "$COMPOSE_FILE" run --rm migrator

echo "Starting API..."
docker compose -f "$COMPOSE_FILE" up -d api

echo "Deployment complete. Checking API health..."
docker compose -f "$COMPOSE_FILE" ps

echo "Done."
