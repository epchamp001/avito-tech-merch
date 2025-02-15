#!/bin/sh

set -e

echo "Ожидание PostgreSQL..."
until pg_isready -h db -U champ001 -d merge-store; do
  sleep 2
done

echo "Применение миграций..."
goose -dir /app/migrations postgres "postgres://champ001:${DB_PASSWORD}@db:5432/merge-store?sslmode=disable" up

echo "Запуск приложения..."
exec /app/merch-store
