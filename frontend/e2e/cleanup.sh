#!/usr/bin/env bash
set -e

# Config (match values used by setup.sh)
DB_NAME="bikemap_test"
DB_USER="testuser"
DB_PASS="testpass"
DB_PORT=5435
PG_CONTAINER_NAME="bikemap_test_db"
BACKEND_BINARY="./tmp/backend"

echo "Stopping backend (if running)..."
if pgrep -f "$BACKEND_BINARY" >/dev/null 2>&1; then
  pkill -f "$BACKEND_BINARY" || true
  sleep 1
fi

echo "Removing backend binary if present..."
if [ -f "$BACKEND_BINARY" ]; then
  rm -f "$BACKEND_BINARY"
fi

echo "Stopping and removing Postgres container ($PG_CONTAINER_NAME)..."
docker rm -f $PG_CONTAINER_NAME 2>/dev/null || true

echo "Unsetting test environment variables..."
unset POSTGRES_TEST_USER POSTGRES_TEST_PASSWORD POSTGRES_TEST_DB POSTGRES_TEST_PORT DATABASE_URL

echo "Cleanup complete."
