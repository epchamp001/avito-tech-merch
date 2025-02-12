package service

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage"
	"context"
	"errors"
	"github.com/google/uuid"
)

type transactionService struct {
	userRepo storage.UserRepository
	txRepo   storage.TransactionRepository
}

func NewTransactionService(userRepo storage.UserRepository, txRepo storage.TransactionRepository) TransactionService {
	return &transactionService{userRepo: userRepo, txRepo: txRepo}
}

func (s *transactionService) TransferCoins(ctx context.Context, senderID, receiverID uuid.UUID, amount int) error {
	if amount <= 0 {
		return errors.New("сумма перевода должна быть положительной")
	}

	sender, err := s.userRepo.GetUserByID(ctx, senderID)
	if err != nil || sender == nil {
		return errors.New("отправитель не найден")
	}

	receiver, err := s.userRepo.GetUserByID(ctx, receiverID)
	if err != nil || receiver == nil {
		return errors.New("получатель не найден")
	}

	if sender.Balance < amount {
		return errors.New("недостаточно монет")
	}

	if err := s.userRepo.UpdateBalance(ctx, senderID, sender.Balance-amount); err != nil {
		return err
	}
	if err := s.userRepo.UpdateBalance(ctx, receiverID, receiver.Balance+amount); err != nil {
		return err
	}

	tx := &models.Transaction{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     amount,
	}
	return s.txRepo.CreateTransaction(ctx, tx)
}

func (s *transactionService) GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	return s.txRepo.GetTransactionsByUser(ctx, userID)
}
