package mocks

import (
	"avito-tech-merch/internal/models"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
	UserRepository
	MerchRepository
	TransactionRepository
	PurchaseRepository
}

// Реализация методов UserRepository
func (m *Repository) CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *Repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *Repository) UpdateBalance(ctx context.Context, userID uuid.UUID, newBalance int) error {
	args := m.Called(ctx, userID, newBalance)
	return args.Error(0)
}

// Реализация методов MerchRepository
func (m *Repository) GetAllMerch(ctx context.Context) ([]models.Merch, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Merch), args.Error(1)
}

func (m *Repository) GetMerchByID(ctx context.Context, merchID int) (*models.Merch, error) {
	args := m.Called(ctx, merchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Merch), args.Error(1)
}

func (m *Repository) GetMerchByName(ctx context.Context, name string) (*models.Merch, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Merch), args.Error(1)
}

// Реализация методов TransactionRepository
func (m *Repository) CreateTransaction(ctx context.Context, tx *models.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *Repository) GetTransactionsByUser(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

// Реализация методов PurchaseRepository
func (m *Repository) CreatePurchase(ctx context.Context, purchase *models.Purchase) error {
	args := m.Called(ctx, purchase)
	return args.Error(0)
}

func (m *Repository) GetPurchasesByUser(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Purchase), args.Error(1)
}
