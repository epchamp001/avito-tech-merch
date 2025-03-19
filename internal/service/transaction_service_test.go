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

func TestTransferCoins_InvalidAmount(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	// Для случая, когда ошибка происходит до вызова транзакции, txManager не используется
	service := NewTransactionService(repoMock, loggerMock, nil)

	ctx := context.Background()
	senderID, receiverID := 1, 2
	amount := 0 // неверная сумма

	loggerMock.
		On("Warnw",
			"Invalid transfer amount",
			"senderID", senderID,
			"receiverID", receiverID,
			"amount", amount,
		).Return()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid transfer amount")
	loggerMock.AssertCalled(t, "Warnw",
		"Invalid transfer amount",
		"senderID", senderID,
		"receiverID", receiverID,
		"amount", amount,
	)
}

func TestTransferCoins_SameSenderReceiver(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	service := NewTransactionService(repoMock, loggerMock, nil)

	ctx := context.Background()
	senderID := 1
	receiverID := 1
	amount := 100

	loggerMock.
		On("Warnw",
			"Sender and receiver are the same",
			"senderID", senderID,
			"receiverID", receiverID,
		).Return()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot transfer to yourself")
	loggerMock.AssertCalled(t, "Warnw",
		"Sender and receiver are the same",
		"senderID", senderID,
		"receiverID", receiverID,
	)
}

func TestTransferCoins_GetSenderBalanceError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewTransactionService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100
	expectedErr := errors.New("sender balance error")

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("failed to get sender balance: %w", expectedErr)).Once()

	repoMock.
		On("GetBalanceByID", mock.Anything, senderID).
		Return(0, expectedErr).Once()

	loggerMock.
		On("Errorw",
			"Failed to get sender balance",
			"senderID", senderID,
			"error", expectedErr,
		).Return().Once()

	loggerMock.
		On("Errorw",
			"Non-retryable error during TransferCoins",
			"error", fmt.Errorf("failed to get sender balance: %w", expectedErr),
		).Return().Once()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get sender balance")

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get sender balance",
		"senderID", senderID,
		"error", expectedErr,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Non-retryable error during TransferCoins",
		"error", fmt.Errorf("failed to get sender balance: %w", expectedErr),
	)
	repoMock.AssertCalled(t, "GetBalanceByID", mock.Anything, senderID)
	txManagerMock.AssertExpectations(t)
}

func TestTransferCoins_InsufficientFunds(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewTransactionService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100
	balance := 50

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("insufficient funds")).Once()

	repoMock.
		On("GetBalanceByID", mock.Anything, senderID).Return(balance, nil).Once()

	loggerMock.
		On("Warnw",
			"Insufficient funds",
			"senderID", senderID,
			"balance", balance,
			"amount", amount,
		).Return().Once()

	loggerMock.
		On("Errorw",
			"Non-retryable error during TransferCoins",
			"error", fmt.Errorf("insufficient funds"),
		).Return().Once()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient funds")

	loggerMock.AssertCalled(t, "Warnw",
		"Insufficient funds",
		"senderID", senderID,
		"balance", balance,
		"amount", amount,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Non-retryable error during TransferCoins",
		"error", fmt.Errorf("insufficient funds"),
	)
	repoMock.AssertCalled(t, "GetBalanceByID", mock.Anything, senderID)
	txManagerMock.AssertExpectations(t)
}

func TestTransferCoins_GetReceiverBalanceError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewTransactionService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100
	senderBalance := 200
	expectedErr := errors.New("receiver balance error")

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("failed to get receiver balance: %w", expectedErr)).Once()

	repoMock.
		On("GetBalanceByID", mock.Anything, senderID).Return(senderBalance, nil).Once()
	repoMock.
		On("GetBalanceByID", mock.Anything, receiverID).Return(0, expectedErr).Once()

	loggerMock.
		On("Errorw",
			"Failed to get receiver balance",
			"receiverID", receiverID,
			"error", expectedErr,
		).Return().Once()
	loggerMock.
		On("Errorw",
			"Non-retryable error during TransferCoins",
			"error", fmt.Errorf("failed to get receiver balance: %w", expectedErr),
		).Return().Once()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get receiver balance")

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get receiver balance",
		"receiverID", receiverID,
		"error", expectedErr,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Non-retryable error during TransferCoins",
		"error", fmt.Errorf("failed to get receiver balance: %w", expectedErr),
	)
	repoMock.AssertCalled(t, "GetBalanceByID", mock.Anything, senderID)
	repoMock.AssertCalled(t, "GetBalanceByID", mock.Anything, receiverID)
	txManagerMock.AssertExpectations(t)
}

func TestTransferCoins_UpdateSenderBalanceError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewTransactionService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100
	senderBalance := 200
	receiverBalance := 100
	expectedErr := errors.New("update sender balance error")
	newSenderBalance := senderBalance - amount

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("failed to update sender balance: %w", expectedErr)).Once()

	repoMock.
		On("GetBalanceByID", mock.Anything, senderID).Return(senderBalance, nil).Once()
	repoMock.
		On("GetBalanceByID", mock.Anything, receiverID).Return(receiverBalance, nil).Once()
	repoMock.
		On("UpdateBalance", mock.Anything, senderID, newSenderBalance).
		Return(expectedErr).Once()

	loggerMock.
		On("Errorw",
			"Failed to update sender balance",
			"senderID", senderID,
			"newBalance", newSenderBalance,
			"error", expectedErr,
		).Return().Once()
	loggerMock.
		On("Errorw",
			"Non-retryable error during TransferCoins",
			"error", fmt.Errorf("failed to update sender balance: %w", expectedErr),
		).Return().Once()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update sender balance")
	repoMock.AssertCalled(t, "UpdateBalance", mock.Anything, senderID, newSenderBalance)
	txManagerMock.AssertExpectations(t)
	loggerMock.AssertCalled(t, "Errorw",
		"Non-retryable error during TransferCoins",
		"error", fmt.Errorf("failed to update sender balance: %w", expectedErr),
	)
}

func TestTransferCoins_UpdateReceiverBalanceError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewTransactionService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100
	senderBalance := 200
	receiverBalance := 50
	newSenderBalance := senderBalance - amount
	newReceiverBalance := receiverBalance + amount
	expectedErr := errors.New("update receiver balance error")

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("failed to update receiver balance: %w", expectedErr)).Once()

	repoMock.
		On("GetBalanceByID", mock.Anything, senderID).Return(senderBalance, nil).Once()
	repoMock.
		On("GetBalanceByID", mock.Anything, receiverID).Return(receiverBalance, nil).Once()
	repoMock.
		On("UpdateBalance", mock.Anything, senderID, newSenderBalance).Return(nil).Once()
	repoMock.
		On("UpdateBalance", mock.Anything, receiverID, newReceiverBalance).Return(expectedErr).Once()

	loggerMock.
		On("Errorw",
			"Failed to update receiver balance",
			"receiverID", receiverID,
			"newBalance", newReceiverBalance,
			"error", expectedErr,
		).Return().Once()

	loggerMock.
		On("Errorw",
			"Non-retryable error during TransferCoins",
			"error", fmt.Errorf("failed to update receiver balance: %w", expectedErr),
		).Return().Once()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update receiver balance")
	repoMock.AssertCalled(t, "UpdateBalance", mock.Anything, receiverID, newReceiverBalance)
	txManagerMock.AssertExpectations(t)
}

func TestTransferCoins_CreateTransactionError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewTransactionService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100
	senderBalance := 200
	receiverBalance := 50
	newSenderBalance := senderBalance - amount
	newReceiverBalance := receiverBalance + amount
	expectedErr := errors.New("create transaction error")

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(fmt.Errorf("failed to create transaction: %w", expectedErr)).Once()

	repoMock.
		On("GetBalanceByID", mock.Anything, senderID).Return(senderBalance, nil).Once()
	repoMock.
		On("GetBalanceByID", mock.Anything, receiverID).Return(receiverBalance, nil).Once()
	repoMock.
		On("UpdateBalance", mock.Anything, senderID, newSenderBalance).Return(nil).Once()
	repoMock.
		On("UpdateBalance", mock.Anything, receiverID, newReceiverBalance).Return(nil).Once()
	repoMock.
		On("CreateTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction")).
		Return(0, expectedErr).Once()

	loggerMock.
		On("Errorw",
			"Failed to create transaction",
			"senderID", senderID,
			"receiverID", receiverID,
			"amount", amount,
			"error", expectedErr,
		).Return().Once()

	loggerMock.
		On("Errorw",
			"Non-retryable error during TransferCoins",
			"error", fmt.Errorf("failed to create transaction: %w", expectedErr),
		).Return().Once()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create transaction")

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to create transaction",
		"senderID", senderID,
		"receiverID", receiverID,
		"amount", amount,
		"error", expectedErr,
	)
	repoMock.AssertCalled(t, "CreateTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction"))
	txManagerMock.AssertExpectations(t)
}

func TestTransferCoins_CommitTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewTransactionService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100
	senderBalance := 200
	receiverBalance := 50
	newSenderBalance := senderBalance - amount
	newReceiverBalance := receiverBalance + amount
	expectedErr := errors.New("commit tx error")

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(expectedErr).Once()

	repoMock.
		On("GetBalanceByID", mock.Anything, senderID).Return(senderBalance, nil).Once()
	repoMock.
		On("GetBalanceByID", mock.Anything, receiverID).Return(receiverBalance, nil).Once()
	repoMock.
		On("UpdateBalance", mock.Anything, senderID, newSenderBalance).Return(nil).Once()
	repoMock.
		On("UpdateBalance", mock.Anything, receiverID, newReceiverBalance).Return(nil).Once()
	repoMock.
		On("CreateTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction")).Return(1, nil).Once()

	loggerMock.
		On("Errorw",
			"Non-retryable error during TransferCoins",
			"error", expectedErr,
		).Return().Once()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "commit tx error")

	loggerMock.AssertCalled(t, "Errorw",
		"Non-retryable error during TransferCoins",
		"error", expectedErr,
	)
	txManagerMock.AssertExpectations(t)
}

func TestTransferCoins_Success(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)
	service := NewTransactionService(repoMock, loggerMock, txManagerMock)

	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100
	senderBalance := 200
	receiverBalance := 50
	newSenderBalance := senderBalance - amount
	newReceiverBalance := receiverBalance + amount

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelSerializable, postgres.AccessModeReadWrite,
			mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(nil)

	repoMock.
		On("GetBalanceByID", mock.Anything, senderID).Return(senderBalance, nil)
	repoMock.
		On("GetBalanceByID", mock.Anything, receiverID).Return(receiverBalance, nil)
	repoMock.
		On("UpdateBalance", mock.Anything, senderID, newSenderBalance).Return(nil)
	repoMock.
		On("UpdateBalance", mock.Anything, receiverID, newReceiverBalance).Return(nil)
	repoMock.
		On("CreateTransaction", mock.Anything, mock.MatchedBy(func(tr *models.Transaction) bool {
			return tr.SenderID == senderID && tr.ReceiverID == receiverID && tr.Amount == amount
		})).Return(1, nil)

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.NoError(t, err)

	repoMock.AssertCalled(t, "GetBalanceByID", mock.Anything, senderID)
	repoMock.AssertCalled(t, "GetBalanceByID", mock.Anything, receiverID)
	repoMock.AssertCalled(t, "UpdateBalance", mock.Anything, senderID, newSenderBalance)
	repoMock.AssertCalled(t, "UpdateBalance", mock.Anything, receiverID, newReceiverBalance)
	repoMock.AssertCalled(t, "CreateTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction"))
	txManagerMock.AssertExpectations(t)
}
