#!/usr/bin/env bash
set -e

# Config
DB_NAME="bikemap_test"
DB_USER="testuser"
DB_PASS="testpass"
DB_PORT=5435
PG_CONTAINER_NAME="bikemap_test_db"
BACKEND_BINARY="./tmp/backend"

# Cleanup function
cleanup() {
    echo "Cleaning up..."
    if [[ -n "$BACKEND_PID" ]]; then
        kill -9 $BACKEND_PID 2>/dev/null || true
    fi
    docker rm -f $PG_CONTAINER_NAME 2>/dev/null || true
    rm -f $BACKEND_BINARY
}
# Ensure cleanup runs on exit or interrupt
trap cleanup EXIT

# Cleanup any existing container
docker rm -f $PG_CONTAINER_NAME 2>/dev/null || true

# Start fresh Postgres container
docker run --rm -d \
  --name $PG_CONTAINER_NAME \
  -e POSTGRES_USER=$DB_USER \
  -e POSTGRES_PASSWORD=$DB_PASS \
  -e POSTGRES_DB=$DB_NAME \
  -p $DB_PORT:5432 \
  postgres:15

# Wait for Postgres to be ready
echo "Waiting for Postgres to start..."
until pg_isready -h localhost -p $DB_PORT -U $DB_USER; do
  sleep 1
done
echo "Postgres is ready!"

# Run migrations
export DATABASE_URL="postgres://$DB_USER:$DB_PASS@localhost:$DB_PORT/$DB_NAME?sslmode=disable"
migrate -path ./migrations -database "$DATABASE_URL" up

# Populate db
echo "Populating test database..."
go run cmd/populate_test_db/main.go

# Build go binary and start the server
go build -o $BACKEND_BINARY cmd/server/main.go
$BACKEND_BINARY &
BACKEND_PID=$!
sleep 2 # give server a moment to start

export POSTGRES_TEST_USER=$DB_USER
export POSTGRES_TEST_PASSWORD=$DB_PASS
export POSTGRES_TEST_DB=$DB_NAME
export POSTGRES_TEST_PORT=$DB_PORT

# Run E2E tests
cd frontend && npx playwright test
