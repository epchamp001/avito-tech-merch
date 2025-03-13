package service

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/db"
	"avito-tech-merch/pkg/logger"
	"context"
	"fmt"
	"time"
)

type transactionService struct {
	repo   db.Repository
	logger logger.Logger
}

func NewTransactionService(repo db.Repository, log logger.Logger) TransactionService {
	return &transactionService{repo: repo, logger: log}
}

func (s *transactionService) TransferCoins(ctx context.Context, senderID int, receiverID int, amount int) error {
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

	var err error
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		s.logger.Errorw("Failed to begin transaction",
			"senderID", senderID,
			"receiverID", receiverID,
			"amount", amount,
			"error", err,
		)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := s.repo.RollbackTx(ctx, tx); rollbackErr != nil {
				s.logger.Errorw("Failed to rollback transaction",
					"senderID", senderID,
					"receiverID", receiverID,
					"amount", amount,
					"error", rollbackErr,
				)
			}
		}
	}()

	senderBalance, err := s.repo.GetBalanceByID(ctx, senderID)
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
		err = fmt.Errorf("insufficient funds")
		return fmt.Errorf("insufficient funds")
	}

	receiverBalance, err := s.repo.GetBalanceByID(ctx, receiverID)
	if err != nil {
		s.logger.Errorw("Failed to get receiver balance",
			"receiverID", receiverID,
			"error", err,
		)
		return fmt.Errorf("failed to get receiver balance: %w", err)
	}

	newSenderBalance := senderBalance - amount
	if err = s.repo.UpdateBalance(ctx, senderID, newSenderBalance); err != nil {
		s.logger.Errorw("Failed to update sender balance",
			"senderID", senderID,
			"newBalance", newSenderBalance,
			"error", err,
		)
		return fmt.Errorf("failed to update sender balance: %w", err)
	}

	newReceiverBalance := receiverBalance + amount
	if err = s.repo.UpdateBalance(ctx, receiverID, newReceiverBalance); err != nil {
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

	if _, err = s.repo.CreateTransaction(ctx, transaction); err != nil {
		s.logger.Errorw("Failed to create transaction",
			"senderID", senderID,
			"receiverID", receiverID,
			"amount", amount,
			"error", err,
		)
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	if err = s.repo.CommitTx(ctx, tx); err != nil {
		s.logger.Errorw("Failed to commit transaction",
			"senderID", senderID,
			"receiverID", receiverID,
			"error", err,
		)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
