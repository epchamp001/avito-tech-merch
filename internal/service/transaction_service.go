package service

import (
	"avito-tech-merch/internal/metrics"
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/internal/storage/db/postgres"
	"avito-tech-merch/pkg/logger"
	"context"
	"fmt"
	"strings"
	"time"
)

type transactionService struct {
	repo      db.Repository
	logger    logger.Logger
	txManager db.TxManager
}

func NewTransactionService(repo db.Repository, log logger.Logger, txManager db.TxManager) TransactionService {
	return &transactionService{repo: repo, logger: log, txManager: txManager}
}

func (s *transactionService) TransferCoins(ctx context.Context, senderID int, receiverID int, amount int) error {
	metrics.RecordCoinTransfer()

	if amount <= 0 {
		s.logger.Warnw("Invalid transfer amount",
			"senderID", senderID,
			"receiverID", receiverID,
			"amount", amount,
		)
		return fmt.Errorf("invalid transfer amount: amount must be positive")
	}

	if senderID == receiverID {
		s.logger.Warnw("Sender and receiver are the same",
			"senderID", senderID,
			"receiverID", receiverID,
		)
		return fmt.Errorf("cannot transfer to yourself")
	}

	const maxRetries = 3
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = s.txManager.WithTx(ctx, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite, func(txCtx context.Context) error {
			senderBalance, err := s.repo.GetBalanceByID(txCtx, senderID)
			if err != nil {
				s.logger.Errorw("Failed to get sender balance",
					"senderID", senderID,
					"error", err,
				)
				return fmt.Errorf("failed to get sender balance: %w", err)
			}

			if senderBalance < amount {
				s.logger.Warnw("Insufficient funds",
					"senderID", senderID,
					"balance", senderBalance,
					"amount", amount,
				)
				return fmt.Errorf("insufficient funds")
			}

			receiverBalance, err := s.repo.GetBalanceByID(txCtx, receiverID)
			if err != nil {
				s.logger.Errorw("Failed to get receiver balance",
					"receiverID", receiverID,
					"error", err,
				)
				return fmt.Errorf("failed to get receiver balance: %w", err)
			}

			newSenderBalance := senderBalance - amount
			if err = s.repo.UpdateBalance(txCtx, senderID, newSenderBalance); err != nil {
				s.logger.Errorw("Failed to update sender balance",
					"senderID", senderID,
					"newBalance", newSenderBalance,
					"error", err,
				)
				return fmt.Errorf("failed to update sender balance: %w", err)
			}

			newReceiverBalance := receiverBalance + amount
			if err = s.repo.UpdateBalance(txCtx, receiverID, newReceiverBalance); err != nil {
				s.logger.Errorw("Failed to update receiver balance",
					"receiverID", receiverID,
					"newBalance", newReceiverBalance,
					"error", err,
				)
				return fmt.Errorf("failed to update receiver balance: %w", err)
			}

			transaction := &models.Transaction{
				SenderID:   senderID,
				ReceiverID: receiverID,
				Amount:     amount,
				CreatedAt:  time.Now(),
			}

			if _, err = s.repo.CreateTransaction(txCtx, transaction); err != nil {
				s.logger.Errorw("Failed to create transaction",
					"senderID", senderID,
					"receiverID", receiverID,
					"amount", amount,
					"error", err,
				)
				return fmt.Errorf("failed to create transaction: %w", err)
			}

			return nil
		})

		if err == nil {
			return nil
		}

		if !IsSerializationError(err) {
			s.logger.Errorw("Non-retryable error during TransferCoins", "error", err)
			return err
		}

		if attempt == maxRetries {
			s.logger.Errorw("Failed to transfer coins after retries", "error", err)
			return err
		}

		s.logger.Infow("Serialization error during coin transfer, retrying", "attempt", attempt, "error", err)
		time.Sleep(time.Duration(attempt) * 100 * time.Millisecond) // Экспоненциальная задержка: 100ms, 200ms, 300ms

	}
	return err
}

func IsSerializationError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "SQLSTATE 40001")
}
