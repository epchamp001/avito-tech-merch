package db

import (
	"avito-tech-merch/internal/models"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresPurchaseRepository struct {
	db *gorm.DB
}

func NewPostgresPurchaseRepository(db *gorm.DB) *postgresPurchaseRepository {
	return &postgresPurchaseRepository{db: db}
}

func (r *postgresPurchaseRepository) CreatePurchase(ctx context.Context, purchase *models.Purchase) error {
	return r.db.WithContext(ctx).Create(purchase).Error
}

func (r *postgresPurchaseRepository) GetPurchasesByUser(ctx context.Context, userID uuid.UUID) ([]models.Purchase, error) {
	var purchases []models.Purchase
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&purchases).Error
	return purchases, err
}
