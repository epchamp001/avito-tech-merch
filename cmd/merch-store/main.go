package main

import (
	_ "avito-tech-merch/docs"
	"avito-tech-merch/internal/app"
	"avito-tech-merch/internal/config"
	"avito-tech-merch/pkg/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Merch Store
// @version 1.0
// @description This is a service that will allow employees to exchange coins and purchase merch with them.

// @contact.name Egor Ponyaev
// @contact.url https://github.com/epchamp001
// @contact.email epchamp001@gmail.com

// @license.name MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token
func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg, err := config.LoadConfig("configs/", ".env")
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(cfg.Env)
	defer log.Sync()

	server := app.NewServer(cfg, log)

	if err := server.Run(); err != nil {
		log.Fatalw("Failed to start server",
			"error", err,
		)
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.PublicServer.ShutdownTimeout)*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Errorw("Shutdown failed",
			"error", err,
		)
		os.Exit(1)
	}

	log.Info("Application stopped gracefully")
}
