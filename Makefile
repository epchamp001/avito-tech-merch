DB_DSN=postgres://champ001:123champ123@localhost:5432/merch-store?sslmode=disable

.PHONY: goose-up goose-down goose-create goose-status goose-reset

# Применить миграции
goose-up:
	goose -dir ./migrations postgres "$(DB_DSN)" up

# Откатить последнюю миграцию
goose-down:
	goose -dir ./migrations postgres "$(DB_DSN)" down

# Посмотреть статус миграций
goose-status:
	goose -dir ./migrations postgres "$(DB_DSN)" status

# Создать новую миграцию
goose-create:
	goose -dir ./migrations create $(name) sql

# Откатить все миграции
goose-reset:
	goose -dir ./migrations postgres "$(DB_DSN)" reset
