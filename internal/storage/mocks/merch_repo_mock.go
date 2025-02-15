package mocks

import (
	"avito-tech-merch/internal/models"
	"context"
	"github.com/stretchr/testify/mock"
)

type MerchRepository struct {
	mock.Mock
}

func (m *MerchRepository) GetAllMerch(ctx context.Context) ([]models.Merch, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Merch), args.Error(1)
}

func (m *MerchRepository) GetMerchByID(ctx context.Context, merchID int) (*models.Merch, error) {
	args := m.Called(ctx, merchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Merch), args.Error(1)
}

func (m *MerchRepository) GetMerchByName(ctx context.Context, name string) (*models.Merch, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Merch), args.Error(1)
}
