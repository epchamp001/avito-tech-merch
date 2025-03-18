package app

import (
	"avito-tech-merch/internal/config"
	controller "avito-tech-merch/internal/controller/http"
	"avito-tech-merch/internal/metrics"
	"avito-tech-merch/internal/service"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/internal/storage/db/postgres"
	"avito-tech-merch/pkg/jwt"
	"avito-tech-merch/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
)

type Server struct {
	closer       *Closer
	router       *gin.Engine
	pgPool       *pgxpool.Pool
	config       *config.Config
	httpServer   *http.Server
	metricServer *http.Server
	logger       logger.Logger
}

func NewServer(cfg *config.Config, log logger.Logger) *Server {
	c := NewCloser()

	pgPool, err := cfg.Storage.ConnectionToPostgres(log)
	if err != nil {
		log.Fatalw("Failed to connect to postgres",
			"error", err)
	}
	c.Add(func(ctx context.Context) error {
		log.Infow("Closing PostgreSQL pool")
		pgPool.Close()
		return nil
	})

	txManager := postgres.NewTxManager(pgPool, log)

	userRepo := postgres.NewUserRepository(txManager, log)
	merchRepo := postgres.NewMerchRepository(txManager, log)
	purchaseRepo := postgres.NewPurchaseRepository(txManager, log)
	transactionRepo := postgres.NewTransactionRepository(txManager, log)

	repo := db.NewRepository(userRepo, merchRepo, purchaseRepo, transactionRepo)

	tokenService := jwt.NewTokenService()

	authService := service.NewAuthService(repo, log, cfg.JWT, tokenService, txManager)
	userService := service.NewUserService(repo, log, txManager)
	merchService := service.NewMerchService(repo, log)
	purchaseService := service.NewPurchaseService(repo, log, txManager)
	transactionService := service.NewTransactionService(repo, log, txManager)

	serv := service.NewService(authService, userService, merchService, purchaseService, transactionService)

	authController := controller.NewAuthController(serv)
	userController := controller.NewUserController(serv)
	merchController := controller.NewMerchController(serv)
	purchaseController := controller.NewPurchaseController(serv)
	transactionController := controller.NewTransactionController(serv)

	contr := controller.NewController(authController, userController, merchController, purchaseController, transactionController)

	router := gin.Default()
	SetupRoutes(router, contr, serv)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.PublicServer.Endpoint, cfg.PublicServer.Port),
		Handler: router,
	}

	metricMux := http.NewServeMux()
	metricMux.Handle("/metrics", metrics.MetricsHandler())
	metricServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Metrics.Endpoint, cfg.Metrics.Port),
		Handler: metricMux,
	}

	return &Server{
		closer:       c,
		router:       router,
		pgPool:       pgPool,
		config:       cfg,
		httpServer:   httpServer,
		metricServer: metricServer,
		logger:       log,
	}
}

func (s *Server) Run(ctx context.Context) error {
	// Сбор метрик активных соединений
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				stats := s.pgPool.Stat()
				metrics.RecordDBActiveConnections(int(stats.TotalConns()))
			case <-ctx.Done():
				s.logger.Infow("Stopping DB metrics collection goroutine")
				return
			}
		}
	}()

	s.closer.Add(func(ctx context.Context) error {
		s.logger.Infow("Shutting down HTTP server")
		return s.httpServer.Shutdown(ctx)
	})

	s.closer.Add(func(ctx context.Context) error {
		s.logger.Infow("Shutting down Metrics server")
		return s.metricServer.Shutdown(ctx)
	})

	go func() {
		s.logger.Infow("Starting HTTP server",
			"address", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatalw("HTTP server error",
				"error", err)
		}
	}()

	go func() {
		s.logger.Infow("Starting Metrics server", "address", s.metricServer.Addr)
		if err := s.metricServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatalw("Metrics server error",
				"error", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.closer.Close(ctx)
}
