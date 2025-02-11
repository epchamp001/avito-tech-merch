package app

import (
	"avito-tech-merch/internal/config"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"os"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func Run(ctx context.Context, cfg *config.Config) {
	logger := setupLogger(cfg.Env)

	router := gin.Default()
	SetupRoutes(router)

	logger.Info("Запуск сервера", "env", cfg.Env, "port", cfg.PublicServer.Port)

	//redisClient := cache.NewRedisClient(ctx, cfg.Redis.Address)

	logger.Info("Подключение к Redis установлено", "address", cfg.Redis.Address)

	address := fmt.Sprintf("%s:%d", cfg.PublicServer.Endpoint, cfg.PublicServer.Port)
	if err := router.Run(address); err != nil {
		logger.Error("Ошибка при запуске сервера", "error", err)
		log.Fatalf("Не удалось запустить сервер: %v\n", err)
	}

	logger.Info("Сервер успешно запущен")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envDev:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
