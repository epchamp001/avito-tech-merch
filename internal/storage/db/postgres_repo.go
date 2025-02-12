package db

import (
	"avito-tech-merch/internal/storage"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	*PostgresUserRepository
	*PostgresMerchRepository
	*PostgresTransactionRepository
	*PostgresPurchaseRepository
}

func NewPostgresRepository(db *gorm.DB) storage.Repository {
	return &PostgresRepository{
		PostgresUserRepository:        NewPostgresUserRepository(db),
		PostgresMerchRepository:       NewPostgresMerchRepository(db),
		PostgresTransactionRepository: NewPostgresTransactionRepository(db),
		PostgresPurchaseRepository:    NewPostgresPurchaseRepository(db),
	}
}
