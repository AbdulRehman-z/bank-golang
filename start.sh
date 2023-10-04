#!/bin/sh

set -e

echo "Migration started"
/app/migrate -path /app/db/migration -database "$DB_URL" -verbose up

echo "Starting server..."
exec "$@"