#!/bin/sh
set -e

echo "Running migrations..."
goose -dir ./migrations postgres "$DB_DSN" up

echo "Starting application..."
exec ./merch-store