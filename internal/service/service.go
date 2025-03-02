package service

import (
	"avito-tech-merch/internal/models"
	"context"
)

type Service interface {
	AuthService
	UserService
	MerchService
	PurchaseService
	TransactionService
}

type AuthService interface {
	Register(ctx context.Context, username string, password string) (string, error)
	Login(ctx context.Context, username string, password string) (string, error)
	ValidateToken(token string) (int, error)
}

type UserService interface {
	GetInfo(ctx context.Context, userID int) (*models.UserInfoResponse, error)
}

type MerchService interface {
	ListMerch(ctx context.Context) ([]*models.Merch, error)
	GetMerch(ctx context.Context, merchID int) (*models.Merch, error)
}

type PurchaseService interface {
	PurchaseMerch(ctx context.Context, userID int, merchName string) error
}

type TransactionService interface {
	TransferCoins(ctx context.Context, senderID int, receiverID int, amount int) error
}

type service struct {
	AuthService
	UserService
	MerchService
	PurchaseService
	TransactionService
}

func NewService(
	auth AuthService,
	user UserService,
	merch MerchService,
	purchase PurchaseService,
	transaction TransactionService,
) Service {
	return &service{
		AuthService:        auth,
		UserService:        user,
		MerchService:       merch,
		PurchaseService:    purchase,
		TransactionService: transaction,
	}
}
