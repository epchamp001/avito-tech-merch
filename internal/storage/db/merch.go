package db

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage"
	"context"
	"gorm.io/gorm"
)

type postgresMerchRepository struct {
	db *gorm.DB
}

func NewPostgresMerchRepository(db *gorm.DB) storage.MerchRepository {
	return &postgresMerchRepository{db: db}
}

func (r *postgresMerchRepository) GetAllMerch(ctx context.Context) ([]models.Merch, error) {
	var merch []models.Merch
	err := r.db.WithContext(ctx).Find(&merch).Error
	return merch, err
}

func (r *postgresMerchRepository) GetMerchByID(ctx context.Context, merchID int) (*models.Merch, error) {
	var merch models.Merch
	err := r.db.WithContext(ctx).First(&merch, "id = ?", merchID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &merch, err
}

func (r *postgresMerchRepository) GetMerchByName(ctx context.Context, name string) (*models.Merch, error) {
	var merch models.Merch
	err := r.db.WithContext(ctx).
		Where("LOWER(name) = LOWER(?)", name).
		First(&merch).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &merch, err
}
