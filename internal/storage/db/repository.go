package db

import (
	"avito-tech-merch/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	UserRepository
	MerchRepository
	PurchaseRepository
	TransactionRepository
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

type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TxManager interface {
	GetExecutor(ctx context.Context) Executor
	WithTx(ctx context.Context, isoLevel pgx.TxIsoLevel, accessMode pgx.TxAccessMode, fn func(ctx context.Context) error) error
}

type postgresRepository struct {
	UserRepository
	MerchRepository
	PurchaseRepository
	TransactionRepository
}

func NewRepository(
	userRepo UserRepository,
	merchRepo MerchRepository,
	purchaseRepo PurchaseRepository,
	transactionRepo TransactionRepository,
) Repository {
	return &postgresRepository{
		UserRepository:        userRepo,
		MerchRepository:       merchRepo,
		PurchaseRepository:    purchaseRepo,
		TransactionRepository: transactionRepo,
	}
}
