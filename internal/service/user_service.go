package service

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage"
	"avito-tech-merch/internal/utils"
	"context"
	"errors"
	"github.com/google/uuid"
)

type userService struct {
	repo storage.Repository
}

func NewUserService(repo storage.Repository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, username, password string) (*models.User, error) {

	existingUser, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("пользователь уже существует")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: hash,
		Balance:      1000, // Начальный баланс
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

func (s *userService) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("ошибка запроса к БД")
	}

	if user == nil {
		user, err = s.RegisterUser(ctx, username, password)
		if err != nil {
			return nil, errors.New("ошибка при регистрации нового пользователя")
		}
	} else {
		if !utils.CheckPasswordHash(password, user.PasswordHash) {
			return nil, errors.New("неверный пароль")
		}
	}

	return user, nil
}

func (s *userService) GetUserInfo(ctx context.Context, userID uuid.UUID) (*models.UserInfoResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("пользователь не найден")
	}

	purchases, err := s.repo.GetPurchasesByUser(ctx, userID)
	if err != nil {
		return nil, errors.New("ошибка получения списка покупок")
	}

	transactions, err := s.repo.GetTransactionsByUser(ctx, userID)
	if err != nil {
		return nil, errors.New("ошибка получения истории транзакций")
	}

	return &models.UserInfoResponse{
		Balance:      user.Balance,
		Purchases:    purchases,
		Transactions: transactions,
	}, nil
}
