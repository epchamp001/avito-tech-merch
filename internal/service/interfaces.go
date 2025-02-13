package service

import (
	"avito-tech-merch/internal/models"
	"context"
	"github.com/google/uuid"
)

type Service interface {
	UserService
	MerchService
	TransactionService
}

type UserService interface {
	RegisterUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateBalance(ctx context.Context, userID uuid.UUID, amount int) error
	AuthenticateUser(ctx context.Context, username, password string) (*models.User, error)
	GetUserInfo(ctx context.Context, userID uuid.UUID) (*models.UserInfoResponse, error)
}

type MerchService interface {
	BuyMerch(ctx context.Context, userID uuid.UUID, merchID int) error
	GetUserPurchases(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error)
	GetMerchByName(ctx context.Context, name string) (*models.Merch, error)
}

type TransactionService interface {
	TransferCoins(ctx context.Context, senderID, receiverID uuid.UUID, amount int) error
	GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error)
}
