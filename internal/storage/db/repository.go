package db

import (
	"avito-tech-merch/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	UserRepository
	MerchRepository
	PurchaseRepository
	TransactionRepository
	TxManager
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (int, error)
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetBalanceByID(ctx context.Context, userID int) (int, error)
	GetBalanceByName(ctx context.Context, username string) (int, error)
	UpdateBalance(ctx context.Context, userID int, newBalance int) error
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

type MerchRepository interface {
	GetAllMerch(ctx context.Context) ([]*models.Merch, error)
	GetMerchByID(ctx context.Context, id int) (*models.Merch, error)
	GetMerchByName(ctx context.Context, merchName string) (*models.Merch, error)
}

type PurchaseRepository interface {
	CreatePurchase(ctx context.Context, purchase *models.Purchase) (int, error)
	GetPurchaseByUserID(ctx context.Context, userID int) ([]*models.Purchase, error)
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *models.Transaction) (int, error)
	GetTransactionByUserID(ctx context.Context, userID int) ([]*models.Transaction, error)
}

type TxManager interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	CommitTx(ctx context.Context, tx pgx.Tx) error
	RollbackTx(ctx context.Context, tx pgx.Tx) error
}

type postgresRepository struct {
	UserRepository
	MerchRepository
	PurchaseRepository
	TransactionRepository
	TxManager
}

func NewRepository(
	userRepo UserRepository,
	merchRepo MerchRepository,
	purchaseRepo PurchaseRepository,
	transactionRepo TransactionRepository,
	txManager TxManager,
) Repository {
	return &postgresRepository{
		UserRepository:        userRepo,
		MerchRepository:       merchRepo,
		PurchaseRepository:    purchaseRepo,
		TransactionRepository: transactionRepo,
		TxManager:             txManager,
	}
}
