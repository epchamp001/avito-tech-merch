package app

import (
	"avito-tech-merch/internal/config"
	"avito-tech-merch/internal/controller/http"
	"avito-tech-merch/internal/service"
	database "avito-tech-merch/internal/storage/db"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	dsn := cfg.Storage.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Ошибка подключения к БД", "error", err)
		log.Fatalf("Ошибка подключения к PostgreSQL: %v\n", err)
	}

	logger.Info("Успешное подключение к PostgreSQL")

	repo := database.NewPostgresRepository(db)

	serv := service.NewService(repo)

	authController := http.NewAuthController(serv, cfg.JWT.SecretKey)
	userController := http.NewUserController(serv)
	merchController := http.NewMerchController(serv)

	router := gin.Default()

	SetupRoutes(router, authController, userController, merchController, cfg.JWT.SecretKey)

	logger.Info("Запуск сервера", "env", cfg.Env, "port", cfg.PublicServer.Port)

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
