package app

import (
	"avito-tech-merch/internal/config"
	controller "avito-tech-merch/internal/controller/http"
	"avito-tech-merch/internal/service"
	database "avito-tech-merch/internal/storage/db"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
	config *config.Config
	server *http.Server
}

func NewServer(cfg *config.Config) *Server {
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

	authController := controller.NewAuthController(serv, cfg.JWT.SecretKey)
	userController := controller.NewUserController(serv)
	merchController := controller.NewMerchController(serv)

	router := gin.Default()
	SetupRoutes(router, authController, userController, merchController, cfg.JWT.SecretKey)

	address := fmt.Sprintf("%s:%d", cfg.PublicServer.Endpoint, cfg.PublicServer.Port)

	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	return &Server{
		router: router,
		db:     db,
		config: cfg,
		server: server,
	}
}

func (s *Server) Run() error {
	log.Println("Запуск сервера на", s.config.PublicServer.Port)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown() error {
	log.Println("Запускаем Graceful Shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Println("Ошибка при остановке сервера:", err)
		return err
	}
	log.Println("Сервер остановлен корректно")

	sqlDB, err := s.db.DB()
	if err == nil {
		sqlDB.Close()
		log.Println("Соединение с БД закрыто")
	}

	return nil
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
