package service

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"context"
	"fmt"
)

type userService struct {
	repo   db.Repository
	logger logger.Logger
}

func NewUserService(repo db.Repository, log logger.Logger) UserService {
	return &userService{repo: repo, logger: log}
}

func (s *userService) GetInfo(ctx context.Context, userID int) (*models.UserInfo, error) {
	var err error

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		s.logger.Errorw("Failed to begin transaction",
			"userID", userID,
			"error", err,
		)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := s.repo.RollbackTx(ctx, tx); rollbackErr != nil {
				s.logger.Errorw("Failed to rollback transaction",
					"userID", userID,
					"error", rollbackErr,
				)
			}
		}
	}()

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

	if err = s.repo.CommitTx(ctx, tx); err != nil {
		s.logger.Errorw("Failed to commit transaction",
			"userID", userID,
			"error", err,
		)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
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
