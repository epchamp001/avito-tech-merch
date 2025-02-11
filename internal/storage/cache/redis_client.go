package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, redisURL string) *RedisClient {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Ошибка парсинга Redis URL: %v", err)
	}

	client := redis.NewClient(opt)
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	return &RedisClient{Client: client}
}
