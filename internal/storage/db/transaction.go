package db

import (
	"avito-tech-merch/internal/models"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresTransactionRepository struct {
	db *gorm.DB
}

func NewPostgresTransactionRepository(db *gorm.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) CreateTransaction(ctx context.Context, tx *models.Transaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *PostgresTransactionRepository) GetTransactionsByUser(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.WithContext(ctx).
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Find(&transactions).Error
	return transactions, err
}
