FROM golang:1.24.1-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o merch-store ./cmd/merch-store

# Финальный образ
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/merch-store .

COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY ./migrations ./migrations

COPY ./configs ./configs

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
