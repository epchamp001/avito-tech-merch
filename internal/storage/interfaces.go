package storage

import (
	"avito-tech-merch/internal/models"
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	UserRepository
	MerchRepository
	TransactionRepository
	PurchaseRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateBalance(ctx context.Context, userID uuid.UUID, newBalance int) error
}

type MerchRepository interface {
	GetAllMerch(ctx context.Context) ([]models.Merch, error)
	GetMerchByID(ctx context.Context, merchID int) (*models.Merch, error)
	GetMerchByName(ctx context.Context, name string) (*models.Merch, error)
}

type PurchaseRepository interface {
	CreatePurchase(ctx context.Context, purchase *models.Purchase) error
	GetPurchasesByUser(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error)
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *models.Transaction) error
	GetTransactionsByUser(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error)
}
