package integration

import (
	"avito-tech-merch/internal/app"
	"avito-tech-merch/internal/config"
	controller "avito-tech-merch/internal/controller/http"
	"avito-tech-merch/internal/service"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/internal/storage/db/postgres"
	"avito-tech-merch/pkg/jwt"
	"avito-tech-merch/pkg/logger"
	"avito-tech-merch/tests/integration/testutil"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
	"time"
)

type TestSuite struct {
	suite.Suite
	psqlContainer *testutil.PostgreSQLContainer
	server        *httptest.Server
}

func (s *TestSuite) SetupSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	cfg, err := config.LoadConfig("../../configs/", "../../.env")
	s.Require().NoError(err)

	psqlContainer, err := testutil.NewPostgreSQLContainer(ctx)
	s.Require().NoError(err)

	s.psqlContainer = psqlContainer

	err = testutil.RunMigrations(psqlContainer.GetDSN(), "../../migrations")
	s.Require().NoError(err)

	poolConfig, err := pgxpool.ParseConfig(psqlContainer.GetDSN())
	s.Require().NoError(err)

	poolConfig.MaxConns = int32(cfg.Storage.Postgres.Pool.MaxConnections)
	poolConfig.MinConns = int32(cfg.Storage.Postgres.Pool.MinConnections)
	poolConfig.MaxConnLifetime = time.Duration(cfg.Storage.Postgres.Pool.MaxLifeTime)
	poolConfig.MaxConnIdleTime = time.Duration(cfg.Storage.Postgres.Pool.MaxIdleTime)
	poolConfig.HealthCheckPeriod = time.Duration(cfg.Storage.Postgres.Pool.HealthCheckPeriod)

	pgPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	s.Require().NoError(err)

	log := logger.NewLogger(cfg.Env)
	defer log.Sync()

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
	app.SetupRoutes(router, contr, serv)

	s.server = httptest.NewServer(router)
}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.psqlContainer.Terminate(ctx))

	s.server.Close()
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
