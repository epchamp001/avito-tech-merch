package app

import (
	"avito-tech-merch/internal/config"
	controller "avito-tech-merch/internal/controller/http"
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
)

type Server struct {
	closer     *Closer
	router     *gin.Engine
	pgPool     *pgxpool.Pool
	config     *config.Config
	httpServer *http.Server
	logger     logger.Logger
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

	userRepo := postgres.NewUserRepository(pgPool, log)
	merchRepo := postgres.NewMerchRepository(pgPool, log)
	purchaseRepo := postgres.NewPurchaseRepository(pgPool, log)
	transactionRepo := postgres.NewTransactionRepository(pgPool, log)
	txManager := postgres.NewTxManager(pgPool, log)

	repo := db.NewRepository(userRepo, merchRepo, purchaseRepo, transactionRepo, txManager)

	tokenService := jwt.NewTokenService()

	authService := service.NewAuthService(repo, log, cfg.JWT, tokenService)
	userService := service.NewUserService(repo, log)
	merchService := service.NewMerchService(repo, log)
	purchaseService := service.NewPurchaseService(repo, log)
	transactionService := service.NewTransactionService(repo, log)

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

	return &Server{
		closer:     c,
		router:     router,
		pgPool:     pgPool,
		config:     cfg,
		httpServer: httpServer,
		logger:     log,
	}
}

func (s *Server) Run() error {
	s.closer.Add(func(ctx context.Context) error {
		s.logger.Infow("Shutting down HTTP server")
		return s.httpServer.Shutdown(ctx)
	})

	go func() {
		s.logger.Infow("Starting HTTP server",
			"address", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatalw("HTTP server error",
				"error", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.closer.Close(ctx)
}
