package mocks

import (
	"avito-tech-merch/internal/models"
	"context"
	"github.com/stretchr/testify/mock"
)

type PurchaseRepository struct {
	mock.Mock
}

func (m *PurchaseRepository) CreatePurchase(ctx context.Context, purchase *models.Purchase) error {
	args := m.Called(ctx, purchase)
	return args.Error(0)
}

func (m *PurchaseRepository) GetPurchasesByUser(ctx context.Context, userID string) ([]models.Purchase, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Purchase), args.Error(1)
}
