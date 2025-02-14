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

# Docker команды
# Собрать Docker-образ
build:
	docker-compose build

# Запустить сервисы
run:
	docker-compose up -d

# Остановить сервисы
stop:
	docker-compose down

# Логи сервиса
logs:
	docker-compose logs -f app

# Применить миграции
migrate:
	$(GOOSE) -dir ./migrations postgres "$(DB_DSN)" up

# Полный сброс БД и миграций
reset-db:
	docker-compose down -v && docker-compose up -d && sleep 5 && $(GOOSE) -dir ./migrations postgres "$(DB_DSN)" up