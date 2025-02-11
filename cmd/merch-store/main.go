package main

import (
	"avito-tech-merch/internal/app"
	"avito-tech-merch/internal/config"
	"context"
	"log"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	cfg, err := config.LoadConfig("configs/config.yaml", ".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v\n", err)
	}

	app.Run(ctx, cfg)
}
