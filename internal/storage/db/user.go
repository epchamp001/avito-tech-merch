package db

import (
	"avito-tech-merch/internal/models"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		First(&user, "id = ?", id).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		First(&user, "username = ?", username).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) UpdateBalance(ctx context.Context, userID uuid.UUID, newBalance int) error {
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("balance", newBalance).
		Error
	return err
}
