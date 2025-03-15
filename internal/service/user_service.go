package service

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/internal/storage/db/postgres"
	"avito-tech-merch/pkg/logger"
	"context"
)

type userService struct {
	repo      db.Repository
	logger    logger.Logger
	txManager db.TxManager
}

func NewUserService(repo db.Repository, log logger.Logger, txManager db.TxManager) UserService {
	return &userService{repo: repo, logger: log, txManager: txManager}
}

func (s *userService) GetInfo(ctx context.Context, userID int) (*models.UserInfo, error) {
	var result *models.UserInfo

	err := s.txManager.WithTx(ctx, postgres.IsolationLevelReadCommitted, postgres.AccessModeReadOnly, func(txCtx context.Context) error {
		user, err := s.repo.GetUserByID(txCtx, userID)
		if err != nil {
			s.logger.Errorw("Failed to get user info",
				"userID", userID,
				"error", err,
			)
			return err
		}

		purchases, err := s.repo.GetPurchaseByUserID(txCtx, userID)
		if err != nil {
			s.logger.Errorw("Failed to get purchases",
				"userID", userID,
				"error", err,
			)
			return err
		}

		transactions, err := s.repo.GetTransactionByUserID(txCtx, userID)
		if err != nil {
			s.logger.Errorw("Failed to get transactions",
				"userID", userID,
				"error", err,
			)
			return err
		}

		result = &models.UserInfo{
			UserID:       user.ID,
			Username:     user.Username,
			Balance:      user.Balance,
			Purchases:    purchases,
			Transactions: transactions,
		}
		return nil
	})
	if err != nil {
		s.logger.Errorw("Error getting user info",
			"error", err,
		)
		return nil, err
	}

	return result, nil
}
