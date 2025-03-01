package main

import (
	"avito-tech-merch/internal/config"
	"avito-tech-merch/pkg/logger"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	log := logger.NewLogger()
	defer log.Sync()

	cfg, err := config.LoadConfig("configs/", ".env", log)
	if err != nil {
		log.Error("Failed to load config", zap.Error(err))
	}

	fmt.Printf("%+v\n", cfg)

	pool, err := cfg.Storage.ConnectionToPostgres(log)
	if err != nil {
		log.Fatal("Failed to connect to postgres", zap.Error(err))
	}
	defer pool.Close()
	log.Info("Successfully connected to PostgreSQL!")

}
