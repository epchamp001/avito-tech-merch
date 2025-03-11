package service

import (
	"avito-tech-merch/internal/models"
	mockRepo "avito-tech-merch/internal/storage/db/mock"
	mockLog "avito-tech-merch/pkg/logger/mock"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestPurchaseMerch_GetMerchError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	service := NewPurchaseService(repoMock, loggerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"

	expectedErr := errors.New("db error")
	repoMock.On("GetMerchByName", ctx, merchName).Return(nil, expectedErr)
	loggerMock.On("Errorw", "Failed to get merch",
		"merchName", merchName,
		"error", expectedErr,
	).Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get merch")

	repoMock.AssertCalled(t, "GetMerchByName", ctx, merchName)
	loggerMock.AssertCalled(t, "Errorw", "Failed to get merch",
		"merchName", merchName,
		"error", expectedErr,
	)
}

func TestPurchaseMerch_GetBalanceError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	service := NewPurchaseService(repoMock, loggerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"

	// Предположим, что мерч найден успешно
	merch := &models.Merch{ID: 10, Name: merchName, Price: 150}
	repoMock.On("GetMerchByName", ctx, merchName).Return(merch, nil)

	expectedErr := errors.New("balance error")
	repoMock.On("GetBalanceByID", ctx, userID).Return(0, expectedErr)
	loggerMock.On("Errorw", "Failed to get user balance",
		"userID", userID,
		"error", expectedErr,
	).Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	repoMock.AssertCalled(t, "GetMerchByName", ctx, merchName)
	repoMock.AssertCalled(t, "GetBalanceByID", ctx, userID)
	loggerMock.AssertCalled(t, "Errorw", "Failed to get user balance",
		"userID", userID,
		"error", expectedErr,
	)
}

func TestPurchaseMerch_InsufficientFunds(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	service := NewPurchaseService(repoMock, loggerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"

	merch := &models.Merch{ID: 10, Name: merchName, Price: 150}
	repoMock.On("GetMerchByName", ctx, merchName).Return(merch, nil)
	// Пользователь имеет меньше средств, чем стоит мерч
	repoMock.On("GetBalanceByID", ctx, userID).Return(100, nil)
	loggerMock.On("Warnw", "Insufficient funds",
		"userID", userID,
		"balance", 100,
		"merchPrice", merch.Price,
	).Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient funds")

	repoMock.AssertCalled(t, "GetMerchByName", ctx, merchName)
	repoMock.AssertCalled(t, "GetBalanceByID", ctx, userID)
	loggerMock.AssertCalled(t, "Warnw", "Insufficient funds",
		"userID", userID,
		"balance", 100,
		"merchPrice", merch.Price,
	)
}

func TestPurchaseMerch_BeginTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	service := NewPurchaseService(repoMock, loggerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"

	merch := &models.Merch{ID: 10, Name: merchName, Price: 50}
	repoMock.On("GetMerchByName", ctx, merchName).Return(merch, nil)
	repoMock.On("GetBalanceByID", ctx, userID).Return(100, nil)

	expectedErr := errors.New("begin tx error")
	repoMock.On("BeginTx", ctx).Return(nil, expectedErr)
	loggerMock.On("Errorw", "Failed to begin transaction",
		"userID", userID,
		"merchName", merchName,
		"error", expectedErr,
	).Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to begin transaction")

	repoMock.AssertCalled(t, "BeginTx", ctx)
	loggerMock.AssertCalled(t, "Errorw", "Failed to begin transaction",
		"userID", userID,
		"merchName", merchName,
		"error", expectedErr,
	)
}

func TestPurchaseMerch_UpdateBalanceError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"

	merch := &models.Merch{ID: 10, Name: merchName, Price: 50}
	repoMock.On("GetMerchByName", ctx, merchName).Return(merch, nil)
	repoMock.On("GetBalanceByID", ctx, userID).Return(100, nil)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)

	newBalance := 100 - merch.Price
	expectedErr := errors.New("update balance error")
	repoMock.On("UpdateBalance", ctx, userID, newBalance).Return(expectedErr)
	// При ошибке обновления баланса происходит откат транзакции
	repoMock.On("RollbackTx", ctx, txMock).Return(nil).Once()

	loggerMock.On("Errorw", "Failed to update user balance",
		"userID", userID,
		"newBalance", newBalance,
		"error", expectedErr,
	).Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user balance")

	repoMock.AssertCalled(t, "UpdateBalance", ctx, userID, newBalance)
	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
	loggerMock.AssertCalled(t, "Errorw", "Failed to update user balance",
		"userID", userID,
		"newBalance", newBalance,
		"error", expectedErr,
	)
}

func TestPurchaseMerch_CreatePurchaseError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"

	merch := &models.Merch{ID: 10, Name: merchName, Price: 50}
	repoMock.On("GetMerchByName", ctx, merchName).Return(merch, nil)
	repoMock.On("GetBalanceByID", ctx, userID).Return(100, nil)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)

	newBalance := 100 - merch.Price
	repoMock.On("UpdateBalance", ctx, userID, newBalance).Return(nil)

	expectedErr := errors.New("create purchase error")
	repoMock.On("CreatePurchase", ctx, mock.AnythingOfType("*models.Purchase")).Return(0, expectedErr)
	loggerMock.On("Errorw", "Failed to create purchase",
		"userID", userID,
		"merchName", merchName,
		"error", expectedErr,
	).Once()

	repoMock.On("RollbackTx", ctx, txMock).Return(nil).Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create purchase")
	repoMock.AssertCalled(t, "CreatePurchase", ctx, mock.AnythingOfType("*models.Purchase"))
	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
	loggerMock.AssertCalled(t, "Errorw", "Failed to create purchase",
		"userID", userID,
		"merchName", merchName,
		"error", expectedErr,
	)
}

func TestPurchaseMerch_CommitTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"

	merch := &models.Merch{ID: 10, Name: merchName, Price: 50}
	repoMock.On("GetMerchByName", ctx, merchName).Return(merch, nil)
	repoMock.On("GetBalanceByID", ctx, userID).Return(100, nil)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)

	newBalance := 100 - merch.Price
	repoMock.On("UpdateBalance", ctx, userID, newBalance).Return(nil)
	repoMock.On("CreatePurchase", ctx, mock.AnythingOfType("*models.Purchase")).Return(1, nil)

	expectedErr := errors.New("commit tx error")
	repoMock.On("CommitTx", ctx, txMock).Return(expectedErr)
	// При ошибке фиксации транзакции происходит откат
	repoMock.On("RollbackTx", ctx, txMock).Return(nil).Once()

	loggerMock.On("Errorw", "Failed to commit transaction",
		"userID", userID,
		"merchName", merchName,
		"error", expectedErr,
	).Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to commit transaction")

	repoMock.AssertCalled(t, "CommitTx", ctx, txMock)
	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
	loggerMock.AssertCalled(t, "Errorw", "Failed to commit transaction",
		"userID", userID,
		"merchName", merchName,
		"error", expectedErr,
	)
}

func TestPurchaseMerch_Success(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"

	merch := &models.Merch{ID: 10, Name: merchName, Price: 50}
	senderBalance := 100
	repoMock.On("GetMerchByName", ctx, merchName).Return(merch, nil)
	repoMock.On("GetBalanceByID", ctx, userID).Return(senderBalance, nil)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)

	newBalance := senderBalance - merch.Price
	repoMock.On("UpdateBalance", ctx, userID, newBalance).Return(nil)
	repoMock.On("CreatePurchase", ctx, mock.MatchedBy(func(p *models.Purchase) bool {
		// Проверяем основные поля покупки
		return p.UserID == userID && p.MerchID == merch.ID
	})).Return(1, nil)
	repoMock.On("CommitTx", ctx, txMock).Return(nil)

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.NoError(t, err)

	repoMock.AssertCalled(t, "GetMerchByName", ctx, merchName)
	repoMock.AssertCalled(t, "GetBalanceByID", ctx, userID)
	repoMock.AssertCalled(t, "BeginTx", ctx)
	repoMock.AssertCalled(t, "UpdateBalance", ctx, userID, newBalance)
	repoMock.AssertCalled(t, "CreatePurchase", ctx, mock.AnythingOfType("*models.Purchase"))
	repoMock.AssertCalled(t, "CommitTx", ctx, txMock)
}

func TestPurchaseMerch_RollbackTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)

	service := NewPurchaseService(repoMock, loggerMock)
	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"

	merch := &models.Merch{
		ID:    10,
		Name:  merchName,
		Price: 50,
	}
	repoMock.On("GetMerchByName", ctx, merchName).Return(merch, nil)
	repoMock.On("GetBalanceByID", ctx, userID).Return(100, nil)

	repoMock.On("BeginTx", ctx).Return(txMock, nil)

	newBalance := 100 - merch.Price

	// Симулируем ошибку при обновлении баланса
	expectedUpdateErr := errors.New("update balance error")
	repoMock.On("UpdateBalance", ctx, userID, newBalance).Return(expectedUpdateErr)

	// Симулируем, что rollback тоже возвращает ошибку
	rollbackErr := errors.New("rollback failure")
	repoMock.On("RollbackTx", ctx, txMock).Return(rollbackErr).Once()

	loggerMock.On("Errorw", "Failed to update user balance",
		"userID", userID,
		"newBalance", newBalance,
		"error", expectedUpdateErr,
	).Once()

	loggerMock.On("Errorw", "Failed to rollback transaction",
		"userID", userID,
		"merchName", merchName,
		"error", rollbackErr,
	).Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user balance")

	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
	loggerMock.AssertCalled(t, "Errorw", "Failed to rollback transaction",
		"userID", userID,
		"merchName", merchName,
		"error", rollbackErr,
	)
}
