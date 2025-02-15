# Базовый образ с Golang
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Устанавливаем goose с статической линковкой
RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@latest

# Сборка бинарного файла
RUN go build -o merch-store ./cmd/merch-store

# Финальный образ
FROM debian:bookworm-slim

WORKDIR /app

# Устанавливаем зависимости (включая PostgreSQL-клиент)
RUN apt-get update && apt-get install -y ca-certificates postgresql-client && rm -rf /var/lib/apt/lists/*

# Копируем бинарник приложения
COPY --from=builder /app/merch-store /app/merch-store

# Копируем `goose` из builder-контейнера в финальный контейнер
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Копируем конфиги, миграции и скрипт запуска
COPY configs/config.yaml /app/configs/config.yaml
COPY .env /app/.env
COPY entrypoint.sh /app/entrypoint.sh
COPY migrations /app/migrations

# Даем права на запуск entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Запуск сервера через entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]