package db

import (
	"avito-tech-merch/internal/models"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresPurchaseRepository struct {
	db *gorm.DB
}

func NewPostgresPurchaseRepository(db *gorm.DB) *PostgresPurchaseRepository {
	return &PostgresPurchaseRepository{db: db}
}

func (r *PostgresPurchaseRepository) CreatePurchase(ctx context.Context, purchase *models.Purchase) error {
	return r.db.WithContext(ctx).Create(purchase).Error
}

func (r *PostgresPurchaseRepository) GetPurchasesByUser(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error) {
	var purchases []models.Purchase
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&purchases).Error
	return purchases, err
}
