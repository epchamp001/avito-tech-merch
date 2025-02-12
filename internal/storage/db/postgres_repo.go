package db

import (
	"avito-tech-merch/internal/storage"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	*postgresUserRepository
	*postgresMerchRepository
	*postgresTransactionRepository
	*postgresPurchaseRepository
}

func NewPostgresRepository(db *gorm.DB) storage.Repository {
	return &PostgresRepository{
		postgresUserRepository:        NewPostgresUserRepository(db),
		postgresMerchRepository:       NewPostgresMerchRepository(db),
		postgresTransactionRepository: NewPostgresTransactionRepository(db),
		postgresPurchaseRepository:    NewPostgresPurchaseRepository(db),
	}
}
