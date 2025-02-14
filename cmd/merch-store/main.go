package main

import (
	"avito-tech-merch/internal/app"
	"avito-tech-merch/internal/config"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig("configs/config.yaml", ".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v\n", err)
	}

	server := app.NewServer(cfg)

	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Остановка сервера...")
	if err := server.Shutdown(); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}
	log.Println("Сервер остановлен корректно")
}
