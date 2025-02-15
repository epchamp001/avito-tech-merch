package service

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/mocks"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestTransferCoins_SameUser(t *testing.T) {
	mockRepo := new(mocks.Repository)
	service := NewTransactionService(mockRepo, mockRepo)

	userID := uuid.New()
	amount := 50

	err := service.TransferCoins(context.TODO(), userID, userID, amount)

	assert.Error(t, err)
	assert.Equal(t, "нельзя отправить монеты самому себе", err.Error())
}

func TestTransferCoins_ReceiverNotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)
	service := NewTransactionService(mockRepo, mockRepo)

	senderID := uuid.New()
	receiverID := uuid.New()
	amount := 50

	mockRepo.On("GetUserByID", mock.Anything, senderID).Return(&models.User{ID: senderID, Balance: 100}, nil)
	mockRepo.On("GetUserByID", mock.Anything, receiverID).Return(nil, errors.New("пользователь не найден"))

	err := service.TransferCoins(context.TODO(), senderID, receiverID, amount)

	assert.Error(t, err)
	assert.Equal(t, "получатель не найден", err.Error())
}

func TestGetUserTransactions(t *testing.T) {
	mockRepo := new(mocks.Repository)
	service := NewTransactionService(mockRepo, mockRepo)

	userID := uuid.New()
	mockRepo.On("GetTransactionsByUser", mock.Anything, userID).Return([]models.Transaction{}, nil)

	transactions, err := service.GetUserTransactions(context.TODO(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, transactions)
}
