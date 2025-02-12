package db

import (
	"avito-tech-merch/internal/models"
	"context"
	"gorm.io/gorm"
)

type PostgresMerchRepository struct {
	db *gorm.DB
}

func NewPostgresMerchRepository(db *gorm.DB) *PostgresMerchRepository {
	return &PostgresMerchRepository{db: db}
}

func (r *PostgresMerchRepository) GetAllMerch(ctx context.Context) ([]models.Merch, error) {
	var merch []models.Merch
	err := r.db.WithContext(ctx).Find(&merch).Error
	return merch, err
}

func (r *PostgresMerchRepository) GetMerchByID(ctx context.Context, merchID int) (*models.Merch, error) {
	var merch models.Merch
	err := r.db.WithContext(ctx).First(&merch, "id = ?", merchID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &merch, err
}

func (r *PostgresMerchRepository) GetMerchByName(ctx context.Context, name string) (*models.Merch, error) {
	var merch models.Merch
	err := r.db.WithContext(ctx).First(&merch, "name = ?", name).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &merch, err
}
