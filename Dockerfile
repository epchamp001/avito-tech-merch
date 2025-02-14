# Используем стабильную версию Go
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git

# Отключаем прокси и проверку контрольных сумм
ENV GOPROXY=direct \
    GOSUMDB=off

# Копируем файлы для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod tidy
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем бинарник
RUN go build -o app ./cmd/merch-store/main.go

# Финальный образ
FROM debian:bullseye-slim

WORKDIR /root/

# Копируем собранное приложение
COPY --from=builder /app/app .
COPY --from=builder /app/configs ./configs

# Запускаем приложение
CMD ["./app"]
