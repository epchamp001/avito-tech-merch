GOOSE=goose
DB_DSN=postgres://champ001:123champ123@localhost:5432/merge-store?sslmode=disable

.PHONY: goose-up goose-down goose-create goose-status goose-reset

# Применить миграции
goose-up:
	$(GOOSE) -dir ./migrations postgres "$(DB_DSN)" up

# Откатить последнюю миграцию
goose-down:
	$(GOOSE) -dir ./migrations postgres "$(DB_DSN)" down

# Посмотреть статус миграций
goose-status:
	$(GOOSE) -dir ./migrations postgres "$(DB_DSN)" status


# Создать новую миграцию
goose-create:
	$(GOOSE) -dir ./migrations create $(name) sql

# Откатить все миграции и применить заново
goose-reset:
	$(GOOSE) -dir ./migrations postgres "$(DB_DSN)" reset

# Собрать Docker-контейнеры без запуска
build:
	docker-compose build --no-cache

# Запуск приложения в Docker
docker-up:
	docker-compose up -d

# Остановка и удаление контейнеров
docker-down:
	docker-compose down

# Пересборка и запуск контейнеров
docker-rebuild:
	docker-compose down
	docker-compose build
	docker-compose up -d

# Просмотр логов приложения
docker-logs:
	docker-compose logs -f app

# Просмотр логов базы данных
docker-logs-db:
	docker-compose logs -f db