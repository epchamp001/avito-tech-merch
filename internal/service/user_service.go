package service

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage"
	"context"
	"errors"
	"github.com/google/uuid"
)

type userService struct {
	repo storage.UserRepository
}

func NewUserService(repo storage.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, username string) (*models.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	// Новый пользователь получает 1000 монет
	newUser := &models.User{
		ID:       uuid.New(),
		Username: username,
		Balance:  1000,
	}

	userID, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	newUser.ID = userID

	return newUser, nil
}

func (s *userService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.repo.GetUserByUsername(ctx, username)
}

func (s *userService) UpdateBalance(ctx context.Context, userID uuid.UUID, amount int) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("пользователь не найден")
	}

	if user.Balance+amount < 0 {
		return errors.New("недостаточно монет")
	}

	return s.repo.UpdateBalance(ctx, userID, user.Balance+amount)
}
