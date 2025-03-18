package service

import (
	"avito-tech-merch/internal/metrics"
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/internal/storage/db/postgres"
	"avito-tech-merch/pkg/logger"
	"context"
	"fmt"
	"time"
)

type purchaseService struct {
	repo      db.Repository
	logger    logger.Logger
	txManager db.TxManager
}

func NewPurchaseService(repo db.Repository, log logger.Logger, txManager db.TxManager) *purchaseService {
	return &purchaseService{repo: repo, logger: log, txManager: txManager}
}

func (s *purchaseService) PurchaseMerch(ctx context.Context, userID int, merchName string) error {
	metrics.RecordMerchPurchase()

	err := s.txManager.WithTx(ctx, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite, func(txCtx context.Context) error {
		merch, err := s.repo.GetMerchByName(txCtx, merchName)
		if err != nil {
			s.logger.Errorw("Failed to get merch",
				"merchName", merchName,
				"error", err,
			)
			return fmt.Errorf("failed to get merch: %w", err)
		}

		balance, err := s.repo.GetBalanceByID(txCtx, userID)
		if err != nil {
			s.logger.Errorw("Failed to get user balance",
				"userID", userID,
				"error", err,
			)
			return fmt.Errorf("failed to get user balance: %w", err)
		}

		if balance < merch.Price {
			s.logger.Warnw("Insufficient funds",
				"userID", userID,
				"balance", balance,
				"merchPrice", merch.Price,
			)
			return fmt.Errorf("insufficient funds")
		}

		newBalance := balance - merch.Price
		if err = s.repo.UpdateBalance(txCtx, userID, newBalance); err != nil {
			s.logger.Errorw("Failed to update user balance",
				"userID", userID,
				"newBalance", newBalance,
				"error", err,
			)
			return fmt.Errorf("failed to update user balance: %w", err)
		}

		purchase := &models.Purchase{
			UserID:    userID,
			MerchID:   merch.ID,
			CreatedAt: time.Now(),
		}

		if _, err = s.repo.CreatePurchase(txCtx, purchase); err != nil {
			s.logger.Errorw("Failed to create purchase",
				"userID", userID,
				"merchName", merchName,
				"error", err,
			)
			return fmt.Errorf("failed to create purchase: %w", err)
		}

		return nil
	})

	if err != nil {
		s.logger.Errorw("Error during PurchaseMerch operation",
			"error", err,
		)
		return err
	}

	return nil
}
