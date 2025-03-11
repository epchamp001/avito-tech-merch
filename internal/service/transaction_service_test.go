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

func TestTransferCoins_InvalidAmount(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	service := NewTransactionService(repoMock, loggerMock)

	ctx := context.Background()
	senderID, receiverID := 1, 2
	amount := 0 // неверная сумма

	loggerMock.On("Warnw",
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
	service := NewTransactionService(repoMock, loggerMock)

	ctx := context.Background()
	senderID := 1
	receiverID := 1
	amount := 100

	loggerMock.On("Warnw",
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
	service := NewTransactionService(repoMock, loggerMock)

	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100

	expectedErr := errors.New("sender balance error")
	repoMock.On("GetBalanceByID", ctx, senderID).Return(0, expectedErr)
	loggerMock.On("Errorw",
		"Failed to get sender balance",
		"senderID", senderID,
		"error", expectedErr,
	).Return()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get sender balance")
	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get sender balance",
		"senderID", senderID,
		"error", expectedErr,
	)
}

func TestTransferCoins_InsufficientFunds(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	service := NewTransactionService(repoMock, loggerMock)

	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100

	// sender balance меньше, чем сумма перевода
	repoMock.On("GetBalanceByID", ctx, senderID).Return(50, nil)
	loggerMock.On("Warnw",
		"Insufficient funds",
		"senderID", senderID,
		"balance", 50,
		"amount", amount,
	).Return()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient funds")
	loggerMock.AssertCalled(t, "Warnw",
		"Insufficient funds",
		"senderID", senderID,
		"balance", 50,
		"amount", amount,
	)
}

func TestTransferCoins_GetReceiverBalanceError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	service := NewTransactionService(repoMock, loggerMock)

	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100

	repoMock.On("GetBalanceByID", ctx, senderID).Return(200, nil)
	// Ошибка получения баланса получателя
	expectedErr := errors.New("receiver balance error")
	repoMock.On("GetBalanceByID", ctx, receiverID).Return(0, expectedErr)
	loggerMock.On("Errorw",
		"Failed to get receiver balance",
		"receiverID", receiverID,
		"error", expectedErr,
	).Return()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get receiver balance")
	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get receiver balance",
		"receiverID", receiverID,
		"error", expectedErr,
	)
}

func TestTransferCoins_BeginTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	service := NewTransactionService(repoMock, loggerMock)

	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100

	// Балансы получены успешно
	repoMock.On("GetBalanceByID", ctx, senderID).Return(200, nil)
	repoMock.On("GetBalanceByID", ctx, receiverID).Return(100, nil)

	expectedErr := errors.New("begin tx error")
	repoMock.On("BeginTx", ctx).Return(nil, expectedErr)
	loggerMock.On("Errorw",
		"Failed to begin transaction",
		"senderID", senderID,
		"receiverID", receiverID,
		"amount", amount,
		"error", expectedErr,
	).Return()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to begin transaction")
	loggerMock.AssertCalled(t, "Errorw",
		"Failed to begin transaction",
		"senderID", senderID,
		"receiverID", receiverID,
		"amount", amount,
		"error", expectedErr,
	)
}

func TestTransferCoins_UpdateSenderBalanceError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)

	service := NewTransactionService(repoMock, loggerMock)
	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100

	repoMock.On("GetBalanceByID", ctx, senderID).Return(200, nil)
	repoMock.On("GetBalanceByID", ctx, receiverID).Return(100, nil)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)

	repoMock.On("RollbackTx", ctx, txMock).Return(nil)

	newSenderBalance := 200 - amount
	expectedErr := errors.New("update sender balance error")
	repoMock.On("UpdateBalance", ctx, senderID, newSenderBalance).Return(expectedErr)

	loggerMock.On("Errorw", "Failed to update sender balance",
		"senderID", senderID,
		"newBalance", newSenderBalance,
		"error", expectedErr,
	)

	err := service.TransferCoins(ctx, senderID, receiverID, amount)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update sender balance")

	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
}

func TestTransferCoins_UpdateReceiverBalanceError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)

	service := NewTransactionService(repoMock, loggerMock)
	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100

	repoMock.On("GetBalanceByID", ctx, senderID).Return(200, nil)
	repoMock.On("GetBalanceByID", ctx, receiverID).Return(50, nil)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)

	newSenderBalance := 200 - amount
	repoMock.On("UpdateBalance", ctx, senderID, newSenderBalance).Return(nil)
	// При обновлении баланса получателя возникает ошибка
	newReceiverBalance := 50 + amount
	expectedErr := errors.New("update receiver balance error")
	repoMock.On("UpdateBalance", ctx, receiverID, newReceiverBalance).Return(expectedErr)
	// Ожидаем, что при ошибке произойдёт откат транзакции
	repoMock.On("RollbackTx", ctx, txMock).Return(nil)

	loggerMock.On("Errorw", "Failed to update receiver balance",
		"receiverID", receiverID,
		"newBalance", newReceiverBalance,
		"error", expectedErr,
	).Once()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update receiver balance")

	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
}

func TestTransferCoins_CreateTransactionError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)

	service := NewTransactionService(repoMock, loggerMock)
	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100

	repoMock.On("GetBalanceByID", ctx, senderID).Return(200, nil)
	repoMock.On("GetBalanceByID", ctx, receiverID).Return(50, nil)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)

	newSenderBalance := 200 - amount
	newReceiverBalance := 50 + amount
	repoMock.On("UpdateBalance", ctx, senderID, newSenderBalance).Return(nil)
	repoMock.On("UpdateBalance", ctx, receiverID, newReceiverBalance).Return(nil)

	// При создании транзакции возникает ошибка
	expectedErr := errors.New("create transaction error")
	repoMock.On("CreateTransaction", ctx, mock.AnythingOfType("*models.Transaction")).Return(0, expectedErr)
	loggerMock.On("Errorw", "Failed to create transaction",
		"senderID", senderID,
		"receiverID", receiverID,
		"amount", amount,
		"error", expectedErr,
	).Once()

	repoMock.On("RollbackTx", ctx, txMock).Return(nil)

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create transaction")

	loggerMock.AssertCalled(t, "Errorw", "Failed to create transaction",
		"senderID", senderID,
		"receiverID", receiverID,
		"amount", amount,
		"error", expectedErr,
	)
	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
}

func TestTransferCoins_CommitTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)

	service := NewTransactionService(repoMock, loggerMock)
	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100

	repoMock.On("GetBalanceByID", ctx, senderID).Return(200, nil)
	repoMock.On("GetBalanceByID", ctx, receiverID).Return(50, nil)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)

	newSenderBalance := 200 - amount
	repoMock.On("UpdateBalance", ctx, senderID, newSenderBalance).Return(nil)
	newReceiverBalance := 50 + amount
	repoMock.On("UpdateBalance", ctx, receiverID, newReceiverBalance).Return(nil)

	repoMock.On("CreateTransaction", ctx, mock.AnythingOfType("*models.Transaction")).Return(1, nil)
	// Ошибка фиксации транзакции
	expectedErr := errors.New("commit tx error")
	repoMock.On("CommitTx", ctx, txMock).Return(expectedErr)
	repoMock.On("RollbackTx", ctx, txMock).Return(nil).Once()

	loggerMock.On("Errorw", "Failed to commit transaction",
		"senderID", senderID,
		"receiverID", receiverID,
		"error", expectedErr,
	).Once()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to commit transaction")
	loggerMock.AssertCalled(t, "Errorw", "Failed to commit transaction",
		"senderID", senderID,
		"receiverID", receiverID,
		"error", expectedErr,
	)
	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)
}

func TestTransferCoins_Success(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)

	service := NewTransactionService(repoMock, loggerMock)
	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100

	senderBalance := 200
	receiverBalance := 50
	repoMock.On("GetBalanceByID", ctx, senderID).Return(senderBalance, nil)
	repoMock.On("GetBalanceByID", ctx, receiverID).Return(receiverBalance, nil)

	repoMock.On("BeginTx", ctx).Return(txMock, nil)

	newSenderBalance := senderBalance - amount
	newReceiverBalance := receiverBalance + amount
	repoMock.On("UpdateBalance", ctx, senderID, newSenderBalance).Return(nil)
	repoMock.On("UpdateBalance", ctx, receiverID, newReceiverBalance).Return(nil)

	repoMock.On("CreateTransaction", ctx, mock.MatchedBy(func(tr *models.Transaction) bool {
		return tr.SenderID == senderID && tr.ReceiverID == receiverID && tr.Amount == amount
	})).Return(1, nil)

	repoMock.On("CommitTx", ctx, txMock).Return(nil)

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.NoError(t, err)

	repoMock.AssertCalled(t, "GetBalanceByID", ctx, senderID)
	repoMock.AssertCalled(t, "GetBalanceByID", ctx, receiverID)
	repoMock.AssertCalled(t, "BeginTx", ctx)
	repoMock.AssertCalled(t, "UpdateBalance", ctx, senderID, newSenderBalance)
	repoMock.AssertCalled(t, "UpdateBalance", ctx, receiverID, newReceiverBalance)
	repoMock.AssertCalled(t, "CreateTransaction", ctx, mock.AnythingOfType("*models.Transaction"))
	repoMock.AssertCalled(t, "CommitTx", ctx, txMock)
}

func TestTransferCoins_RollbackTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)

	service := NewTransactionService(repoMock, loggerMock)
	ctx := context.Background()
	senderID, receiverID, amount := 1, 2, 100

	repoMock.On("GetBalanceByID", ctx, senderID).Return(200, nil)
	repoMock.On("GetBalanceByID", ctx, receiverID).Return(50, nil)
	repoMock.On("BeginTx", ctx).Return(txMock, nil)

	// Ошибка обновления баланса отправителя
	newSenderBalance := 200 - amount
	expectedErr := errors.New("update sender balance error")
	repoMock.On("UpdateBalance", ctx, senderID, newSenderBalance).Return(expectedErr)
	// Обновление баланса получателя не производится

	// Симулируем ошибку при откате транзакции
	rollbackErr := errors.New("rollback failure")
	repoMock.On("RollbackTx", ctx, txMock).Return(rollbackErr).Once()

	loggerMock.On("Errorw", "Failed to update sender balance",
		"senderID", senderID,
		"newBalance", newSenderBalance,
		"error", expectedErr,
	).Once()

	loggerMock.On("Errorw", "Failed to rollback transaction",
		"senderID", senderID,
		"receiverID", receiverID,
		"amount", amount,
		"error", rollbackErr,
	).Once()

	err := service.TransferCoins(ctx, senderID, receiverID, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update sender balance")

	repoMock.AssertCalled(t, "RollbackTx", ctx, txMock)

	loggerMock.AssertCalled(t, "Errorw", "Failed to update sender balance",
		"senderID", senderID,
		"newBalance", newSenderBalance,
		"error", expectedErr,
	)

	loggerMock.AssertCalled(t, "Errorw", "Failed to rollback transaction",
		"senderID", senderID,
		"receiverID", receiverID,
		"amount", amount,
		"error", rollbackErr,
	)
}
