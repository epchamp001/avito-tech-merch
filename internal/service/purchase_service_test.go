package service

import (
	"avito-tech-merch/internal/models"
	mockRepo "avito-tech-merch/internal/storage/db/mock"
	"avito-tech-merch/internal/storage/db/postgres"
	mockLog "avito-tech-merch/pkg/logger/mock"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestPurchaseMerch_GetMerchError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"
	expectedErr := errors.New("db error")

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(nil, expectedErr).Once()

	loggerMock.
		On("Errorw",
			"Failed to get merch",
			"merchName", merchName,
			"error", expectedErr,
		).Return().Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)

	assert.True(t, errors.Is(err, expectedErr), "expected error to wrap %v but got %v", expectedErr, err)

	repoMock.AssertCalled(t, "GetMerchByName", mock.Anything, merchName)
	txManagerMock.AssertNotCalled(t, "WithTx", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get merch",
		"merchName", merchName,
		"error", expectedErr,
	)
}

func TestPurchaseMerch_GetBalanceError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"
	expectedErr := errors.New("balance error")
	merch := &models.Merch{ID: 10, Price: 150}

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(merch, nil).Once()

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("failed to get user balance: %w", expectedErr)).Once()

	repoMock.
		On("GetBalanceByID", mock.Anything, userID).
		Return(0, expectedErr).Once()

	loggerMock.
		On("Errorw",
			"Failed to get user balance",
			"userID", userID,
			"error", expectedErr,
		).Return().Once()

	loggerMock.
		On("Errorw",
			"Non-retryable error during PurchaseMerch",
			"error", fmt.Errorf("failed to get user balance: %w", expectedErr),
		).Return().Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get user balance")

	repoMock.AssertCalled(t, "GetMerchByName", mock.Anything, merchName)
	repoMock.AssertCalled(t, "GetBalanceByID", mock.Anything, userID)
	txManagerMock.AssertExpectations(t)
	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get user balance",
		"userID", userID,
		"error", expectedErr,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Non-retryable error during PurchaseMerch",
		"error", fmt.Errorf("failed to get user balance: %w", expectedErr),
	)
}

func TestPurchaseMerch_InsufficientFunds(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"
	merch := &models.Merch{ID: 10, Name: merchName, Price: 150}

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(merch, nil).Once()

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("insufficient funds")).Once()

	repoMock.
		On("GetBalanceByID", mock.Anything, userID).
		Return(100, nil).Once()

	loggerMock.
		On("Warnw",
			"Insufficient funds",
			"userID", userID,
			"balance", 100,
			"merchPrice", merch.Price,
		).Return().Once()

	loggerMock.
		On("Errorw",
			"Non-retryable error during PurchaseMerch",
			"error", fmt.Errorf("insufficient funds"),
		).Return().Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient funds")

	repoMock.AssertCalled(t, "GetMerchByName", mock.Anything, merchName)
	repoMock.AssertCalled(t, "GetBalanceByID", mock.Anything, userID)
	txManagerMock.AssertExpectations(t)
	loggerMock.AssertCalled(t, "Warnw",
		"Insufficient funds",
		"userID", userID,
		"balance", 100,
		"merchPrice", merch.Price,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Non-retryable error during PurchaseMerch",
		"error", fmt.Errorf("insufficient funds"),
	)
}

func TestPurchaseMerch_BeginTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"
	expectedErr := errors.New("begin tx error")

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(&models.Merch{ID: 10, Price: 50}, nil).Once()

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Return(expectedErr).Once()

	loggerMock.
		On("Errorw",
			"Non-retryable error during PurchaseMerch",
			"error", expectedErr,
		).Return().Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)

	assert.Contains(t, err.Error(), expectedErr.Error())

	repoMock.AssertCalled(t, "GetMerchByName", mock.Anything, merchName)
	txManagerMock.AssertExpectations(t)
	loggerMock.AssertCalled(t, "Errorw",
		"Non-retryable error during PurchaseMerch",
		"error", expectedErr,
	)
}

func TestPurchaseMerch_UpdateBalanceError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"
	merch := &models.Merch{ID: 10, Name: merchName, Price: 50}
	expectedErr := errors.New("update balance error")

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(merch, nil).Once()

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			_ = args.Get(3).(func(context.Context) error)(context.Background())
		}).
		Return(fmt.Errorf("failed to update user balance: %w", expectedErr)).Once()

	repoMock.
		On("GetBalanceByID", mock.Anything, userID).
		Return(100, nil).Once()
	newBalance := 100 - merch.Price
	repoMock.
		On("UpdateBalance", mock.Anything, userID, newBalance).
		Return(expectedErr).Once()

	loggerMock.
		On("Errorw",
			"Failed to update user balance",
			"userID", userID,
			"newBalance", newBalance,
			"error", expectedErr,
		).Return().Once()

	loggerMock.
		On("Errorw",
			"Non-retryable error during PurchaseMerch",
			"error", fmt.Errorf("failed to update user balance: %w", expectedErr),
		).Return().Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user balance")
	assert.True(t, errors.Is(err, expectedErr), "expected error to wrap %v but got %v", expectedErr, err)

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to update user balance",
		"userID", userID,
		"newBalance", newBalance,
		"error", expectedErr,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Non-retryable error during PurchaseMerch",
		"error", fmt.Errorf("failed to update user balance: %w", expectedErr),
	)
	repoMock.AssertCalled(t, "UpdateBalance", mock.Anything, userID, newBalance)
	txManagerMock.AssertExpectations(t)
}

func TestPurchaseMerch_CreatePurchaseError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"
	merch := &models.Merch{ID: 10, Name: merchName, Price: 50}
	expectedErr := errors.New("create purchase error")

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(merch, nil).Once()

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			_ = args.Get(3).(func(context.Context) error)(context.Background())
		}).
		Return(fmt.Errorf("failed to create purchase: %w", expectedErr)).Once()

	repoMock.
		On("GetBalanceByID", mock.Anything, userID).
		Return(100, nil).Once()
	newBalance := 100 - merch.Price
	repoMock.
		On("UpdateBalance", mock.Anything, userID, newBalance).
		Return(nil).Once()
	repoMock.
		On("CreatePurchase", mock.Anything, mock.AnythingOfType("*models.Purchase")).
		Return(0, expectedErr).Once()

	loggerMock.
		On("Errorw",
			"Failed to create purchase",
			"userID", userID,
			"merchName", merchName,
			"error", expectedErr,
		).Return().Once()

	loggerMock.
		On("Errorw",
			"Non-retryable error during PurchaseMerch",
			"error", fmt.Errorf("failed to create purchase: %w", expectedErr),
		).Return().Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create purchase")
	assert.True(t, errors.Is(err, expectedErr), "expected error to wrap %v but got %v", expectedErr, err)

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to create purchase",
		"userID", userID,
		"merchName", merchName,
		"error", expectedErr,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Non-retryable error during PurchaseMerch",
		"error", fmt.Errorf("failed to create purchase: %w", expectedErr),
	)
	repoMock.AssertCalled(t, "CreatePurchase", mock.Anything, mock.AnythingOfType("*models.Purchase"))
	txManagerMock.AssertExpectations(t)
}

func TestPurchaseMerch_Success(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"
	merch := &models.Merch{ID: 10, Name: merchName, Price: 50}
	balance := 100
	newBalance := balance - merch.Price

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(merch, nil).Once()

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			_ = args.Get(3).(func(context.Context) error)(context.Background())
		}).
		Return(nil).Once()

	repoMock.
		On("GetBalanceByID", mock.Anything, userID).
		Return(balance, nil).Once()
	repoMock.
		On("UpdateBalance", mock.Anything, userID, newBalance).
		Return(nil).Once()
	repoMock.
		On("CreatePurchase", mock.Anything, mock.MatchedBy(func(p *models.Purchase) bool {
			return p.UserID == userID && p.MerchID == merch.ID
		})).
		Return(1, nil).Once()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.NoError(t, err)

	repoMock.AssertCalled(t, "GetMerchByName", mock.Anything, merchName)
	repoMock.AssertCalled(t, "GetBalanceByID", mock.Anything, userID)
	repoMock.AssertCalled(t, "UpdateBalance", mock.Anything, userID, newBalance)
	repoMock.AssertCalled(t, "CreatePurchase", mock.Anything, mock.AnythingOfType("*models.Purchase"))
	txManagerMock.AssertExpectations(t)
}
