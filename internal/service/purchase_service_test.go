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

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(expectedErr)

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(nil, expectedErr)

	loggerMock.
		On("Errorw",
			"Failed to get merch",
			"merchName", merchName,
			"error", expectedErr,
		).Return()

	loggerMock.
		On("Errorw",
			"Error during PurchaseMerch operation",
			"error", expectedErr,
		).Return()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get merch",
		"merchName", merchName,
		"error", expectedErr,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Error during PurchaseMerch operation",
		"error", expectedErr,
	)
	repoMock.AssertCalled(t, "GetMerchByName", mock.Anything, merchName)
	txManagerMock.AssertExpectations(t)
}

func TestPurchaseMerch_GetBalanceError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"
	merch := &models.Merch{ID: 10, Name: merchName, Price: 150}
	expectedErr := errors.New("balance error")

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("failed to get user balance: %w", expectedErr))

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(merch, nil)
	repoMock.
		On("GetBalanceByID", mock.Anything, userID).
		Return(0, expectedErr)

	loggerMock.
		On("Errorw",
			"Failed to get user balance",
			"userID", userID,
			"error", expectedErr,
		).Return()
	loggerMock.
		On("Errorw",
			"Error during PurchaseMerch operation",
			"error", fmt.Errorf("failed to get user balance: %w", expectedErr),
		).Return()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get user balance",
		"userID", userID,
		"error", expectedErr,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Error during PurchaseMerch operation",
		"error", fmt.Errorf("failed to get user balance: %w", expectedErr),
	)
	repoMock.AssertCalled(t, "GetMerchByName", mock.Anything, merchName)
	repoMock.AssertCalled(t, "GetBalanceByID", mock.Anything, userID)
	txManagerMock.AssertExpectations(t)
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

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("insufficient funds"))

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(merch, nil)
	repoMock.
		On("GetBalanceByID", mock.Anything, userID).
		Return(100, nil)

	loggerMock.
		On("Warnw",
			"Insufficient funds",
			"userID", userID,
			"balance", 100,
			"merchPrice", merch.Price,
		).Return()
	loggerMock.
		On("Errorw",
			"Error during PurchaseMerch operation",
			"error", fmt.Errorf("insufficient funds"),
		).Return()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient funds")

	loggerMock.AssertCalled(t, "Warnw",
		"Insufficient funds",
		"userID", userID,
		"balance", 100,
		"merchPrice", merch.Price,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Error during PurchaseMerch operation",
		"error", fmt.Errorf("insufficient funds"),
	)
	repoMock.AssertCalled(t, "GetMerchByName", mock.Anything, merchName)
	repoMock.AssertCalled(t, "GetBalanceByID", mock.Anything, userID)
	txManagerMock.AssertExpectations(t)
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

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Return(expectedErr)

	loggerMock.
		On("Errorw",
			"Error during PurchaseMerch operation",
			"error", expectedErr,
		).Return()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())

	loggerMock.AssertCalled(t, "Errorw",
		"Error during PurchaseMerch operation",
		"error", expectedErr,
	)
	txManagerMock.AssertExpectations(t)
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

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("failed to update user balance: %w", expectedErr))

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(merch, nil)
	repoMock.
		On("GetBalanceByID", mock.Anything, userID).
		Return(100, nil)
	newBalance := 100 - merch.Price
	repoMock.
		On("UpdateBalance", mock.Anything, userID, newBalance).
		Return(expectedErr)

	loggerMock.
		On("Errorw",
			"Failed to update user balance",
			"userID", userID,
			"newBalance", newBalance,
			"error", expectedErr,
		).Return()
	loggerMock.
		On("Errorw",
			"Error during PurchaseMerch operation",
			"error", fmt.Errorf("failed to update user balance: %w", expectedErr),
		).Return()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user balance")

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to update user balance",
		"userID", userID,
		"newBalance", newBalance,
		"error", expectedErr,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Error during PurchaseMerch operation",
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

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("failed to create purchase: %w", expectedErr))

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(merch, nil)
	repoMock.
		On("GetBalanceByID", mock.Anything, userID).
		Return(100, nil)
	newBalance := 100 - merch.Price
	repoMock.
		On("UpdateBalance", mock.Anything, userID, newBalance).
		Return(nil)
	repoMock.
		On("CreatePurchase", mock.Anything, mock.AnythingOfType("*models.Purchase")).
		Return(0, expectedErr)

	loggerMock.
		On("Errorw",
			"Failed to create purchase",
			"userID", userID,
			"merchName", merchName,
			"error", expectedErr,
		).Return()
	loggerMock.
		On("Errorw",
			"Error during PurchaseMerch operation",
			"error", fmt.Errorf("failed to create purchase: %w", expectedErr),
		).Return()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create purchase")

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to create purchase",
		"userID", userID,
		"merchName", merchName,
		"error", expectedErr,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Error during PurchaseMerch operation",
		"error", fmt.Errorf("failed to create purchase: %w", expectedErr),
	)
	repoMock.AssertCalled(t, "CreatePurchase", mock.Anything, mock.AnythingOfType("*models.Purchase"))
	txManagerMock.AssertExpectations(t)
}

func TestPurchaseMerch_CommitTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"
	expectedErr := errors.New("commit tx error")

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Return(expectedErr)

	loggerMock.
		On("Errorw",
			"Error during PurchaseMerch operation",
			"error", expectedErr,
		).Return()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commit tx error")

	loggerMock.AssertCalled(t, "Errorw",
		"Error during PurchaseMerch operation",
		"error", expectedErr,
	)
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

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(nil)

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(merch, nil)
	repoMock.
		On("GetBalanceByID", mock.Anything, userID).
		Return(balance, nil)
	repoMock.
		On("UpdateBalance", mock.Anything, userID, newBalance).
		Return(nil)
	repoMock.
		On("CreatePurchase", mock.Anything, mock.MatchedBy(func(p *models.Purchase) bool {
			return p.UserID == userID && p.MerchID == merch.ID
		})).Return(1, nil)

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.NoError(t, err)

	repoMock.AssertCalled(t, "GetMerchByName", mock.Anything, merchName)
	repoMock.AssertCalled(t, "GetBalanceByID", mock.Anything, userID)
	repoMock.AssertCalled(t, "UpdateBalance", mock.Anything, userID, newBalance)
	repoMock.AssertCalled(t, "CreatePurchase", mock.Anything, mock.AnythingOfType("*models.Purchase"))
	txManagerMock.AssertExpectations(t)
}

func TestPurchaseMerch_RollbackTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewPurchaseService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	userID := 1
	merchName := "T-Shirt"
	merch := &models.Merch{ID: 10, Name: merchName, Price: 50}
	balance := 100
	expectedUpdateErr := errors.New("update balance error")
	newBalance := balance - merch.Price

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("failed to update user balance: %w", expectedUpdateErr))

	repoMock.
		On("GetMerchByName", mock.Anything, merchName).
		Return(merch, nil)
	repoMock.
		On("GetBalanceByID", mock.Anything, userID).
		Return(balance, nil)
	repoMock.
		On("UpdateBalance", mock.Anything, userID, newBalance).
		Return(expectedUpdateErr)

	loggerMock.
		On("Errorw",
			"Failed to update user balance",
			"userID", userID,
			"newBalance", newBalance,
			"error", expectedUpdateErr,
		).Return()
	loggerMock.
		On("Errorw",
			"Error during PurchaseMerch operation",
			"error", fmt.Errorf("failed to update user balance: %w", expectedUpdateErr),
		).Return()

	err := service.PurchaseMerch(ctx, userID, merchName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user balance")

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to update user balance",
		"userID", userID,
		"newBalance", newBalance,
		"error", expectedUpdateErr,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Error during PurchaseMerch operation",
		"error", fmt.Errorf("failed to update user balance: %w", expectedUpdateErr),
	)
	txManagerMock.AssertExpectations(t)
}
