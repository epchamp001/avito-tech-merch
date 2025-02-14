FROM golang:1.23 AS builder

WORKDIR /app

ENV GOPROXY=direct
ENV GOSUMDB=off

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/bin/merch-store ./cmd/merch-store/main.go

FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/merch-store /app/merch-store

COPY configs /app/configs

EXPOSE 8080

CMD ["/app/merch-store"]
