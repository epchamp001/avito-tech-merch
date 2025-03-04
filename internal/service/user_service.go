package service

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"context"
)

type userService struct {
	repo   db.Repository
	logger logger.Logger
}

func NewUserService(repo db.Repository, log logger.Logger) UserService {
	return &userService{repo: repo, logger: log}
}

func (s *userService) GetInfo(ctx context.Context, userID int) (*models.UserInfo, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Errorw("Failed to get user info",
			"userID", userID,
			"error", err,
		)
		return nil, err
	}

	purchases, err := s.repo.GetPurchaseByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorw("Failed to get purchases",
			"userID", userID,
			"error", err,
		)
		return nil, err
	}

	transactions, err := s.repo.GetTransactionByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorw("Failed to get transactions",
			"userID", userID,
			"error", err,
		)
		return nil, err
	}

	response := &models.UserInfo{
		UserID:       user.ID,
		Username:     user.Username,
		Balance:      user.Balance,
		Purchases:    purchases,
		Transactions: transactions,
	}

	return response, nil
}
