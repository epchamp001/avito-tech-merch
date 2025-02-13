package db

import (
	"avito-tech-merch/internal/storage"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	storage.UserRepository
	storage.MerchRepository
	storage.TransactionRepository
	storage.PurchaseRepository
}

func NewPostgresRepository(db *gorm.DB) storage.Repository {
	return &PostgresRepository{
		UserRepository:        NewPostgresUserRepository(db),
		MerchRepository:       NewPostgresMerchRepository(db),
		TransactionRepository: NewPostgresTransactionRepository(db),
		PurchaseRepository:    NewPostgresPurchaseRepository(db),
	}
}
