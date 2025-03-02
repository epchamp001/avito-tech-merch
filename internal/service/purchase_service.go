package service

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"context"
	"fmt"
	"time"
)

type purchaseService struct {
	repo   db.Repository
	logger logger.Logger
}

func NewPurchaseService(repo db.Repository, log logger.Logger) *purchaseService {
	return &purchaseService{repo: repo, logger: log}
}

func (s *purchaseService) PurchaseMerch(ctx context.Context, userID int, merchName string) error {
	merch, err := s.repo.GetMerchByName(ctx, merchName)
	if err != nil {
		s.logger.Errorw("Failed to get merch",
			"merchName", merchName,
			"error", err,
		)
		return fmt.Errorf("failed to get merch: %w", err)
	}

	balance, err := s.repo.GetBalanceByID(ctx, userID)
	if err != nil {
		s.logger.Errorw("Failed to get user balance",
			"userID", userID,
			"error", err,
		)
		return err
	}

	if balance < merch.Price {
		s.logger.Warnw("Insufficient funds",
			"userID", userID,
			"balance", balance,
			"merchPrice", merch.Price,
		)
		return fmt.Errorf("insufficient funds")
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		s.logger.Errorw("Failed to begin transaction",
			"userID", userID,
			"merchName", merchName,
			"error", err,
		)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := s.repo.RollbackTx(ctx, tx); rollbackErr != nil {
				s.logger.Errorw("Failed to rollback transaction",
					"userID", userID,
					"merchName", merchName,
					"error", rollbackErr,
				)
			}
		}
	}()

	newBalance := balance - merch.Price
	if err := s.repo.UpdateBalance(ctx, userID, newBalance); err != nil {
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

	if _, err := s.repo.CreatePurchase(ctx, purchase); err != nil {
		s.logger.Errorw("Failed to create purchase",
			"userID", userID,
			"merchName", merchName,
			"error", err,
		)
		return fmt.Errorf("failed to create purchase: %w", err)
	}

	if err := s.repo.CommitTx(ctx, tx); err != nil {
		s.logger.Errorw("Failed to commit transaction",
			"userID", userID,
			"merchName", merchName,
			"error", err,
		)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
