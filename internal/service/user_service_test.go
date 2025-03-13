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
	"time"
)

func TestGetInfo_Success(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)

	expectedUser := &models.User{
		ID:       1,
		Username: "user1",
		Balance:  100,
	}

	expectedPurchases := []*models.Purchase{
		{ID: 1, UserID: 1, MerchID: 10, CreatedAt: time.Now()},
	}

	expectedTransactions := []*models.Transaction{
		{ID: 1, SenderID: 1, ReceiverID: 2, Amount: 50, CreatedAt: time.Now()},
	}

	repoMock.On("BeginTx", mock.Anything).Return(txMock, nil)
	repoMock.On("GetUserByID", mock.Anything, expectedUser.ID).Return(expectedUser, nil)
	repoMock.On("GetPurchaseByUserID", mock.Anything, expectedUser.ID).Return(expectedPurchases, nil)
	repoMock.On("GetTransactionByUserID", mock.Anything, expectedUser.ID).Return(expectedTransactions, nil)
	repoMock.On("CommitTx", mock.Anything, txMock).Return(nil)

	service := NewUserService(repoMock, loggerMock)

	userInfo, err := service.GetInfo(context.Background(), expectedUser.ID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, userInfo.UserID)
	assert.Equal(t, expectedUser.Username, userInfo.Username)
	assert.Equal(t, expectedUser.Balance, userInfo.Balance)
	assert.Equal(t, expectedPurchases, userInfo.Purchases)
	assert.Equal(t, expectedTransactions, userInfo.Transactions)

	repoMock.AssertExpectations(t)
	txMock.AssertExpectations(t)
}

func TestGetInfo_Error_GetUser(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)

	expectedError := assert.AnError
	userID := 1

	repoMock.On("BeginTx", mock.Anything).Return(txMock, nil)
	repoMock.On("GetUserByID", mock.Anything, userID).Return(nil, expectedError)
	repoMock.On("RollbackTx", mock.Anything, txMock).Return(nil)

	loggerMock.On("Errorw",
		"Failed to get user info",
		"userID", userID,
		"error", expectedError,
	).Return()

	service := NewUserService(repoMock, loggerMock)

	userInfo, err := service.GetInfo(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, userInfo)

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get user info",
		"userID", userID,
		"error", expectedError,
	)

	repoMock.AssertCalled(t, "BeginTx", mock.Anything)
	repoMock.AssertCalled(t, "GetUserByID", mock.Anything, userID)
	repoMock.AssertCalled(t, "RollbackTx", mock.Anything, txMock)
}

func TestGetInfo_Error_GetPurchases(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)

	expectedError := assert.AnError
	expectedUser := &models.User{
		ID:       1,
		Username: "user1",
		Balance:  100,
	}

	repoMock.On("BeginTx", mock.Anything).Return(txMock, nil)
	repoMock.On("GetUserByID", mock.Anything, expectedUser.ID).Return(expectedUser, nil)
	repoMock.On("GetPurchaseByUserID", mock.Anything, expectedUser.ID).Return(nil, expectedError)
	repoMock.On("RollbackTx", mock.Anything, txMock).Return(nil)

	loggerMock.On("Errorw",
		"Failed to get purchases",
		"userID", expectedUser.ID,
		"error", expectedError,
	).Return()

	service := NewUserService(repoMock, loggerMock)

	userInfo, err := service.GetInfo(context.Background(), expectedUser.ID)

	assert.Error(t, err)
	assert.Nil(t, userInfo)

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get purchases",
		"userID", expectedUser.ID,
		"error", expectedError,
	)

	repoMock.AssertCalled(t, "BeginTx", mock.Anything)
	repoMock.AssertCalled(t, "GetUserByID", mock.Anything, expectedUser.ID)
	repoMock.AssertCalled(t, "GetPurchaseByUserID", mock.Anything, expectedUser.ID)
	repoMock.AssertCalled(t, "RollbackTx", mock.Anything, txMock)
}

func TestGetInfo_Error_GetTransactions(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)

	expectedError := assert.AnError
	expectedUser := &models.User{
		ID:       1,
		Username: "user1",
		Balance:  100,
	}
	expectedPurchases := []*models.Purchase{
		{ID: 1, UserID: 1, MerchID: 10, CreatedAt: time.Now()},
	}

	repoMock.On("BeginTx", mock.Anything).Return(txMock, nil)
	repoMock.On("GetUserByID", mock.Anything, expectedUser.ID).Return(expectedUser, nil)
	repoMock.On("GetPurchaseByUserID", mock.Anything, expectedUser.ID).Return(expectedPurchases, nil)
	repoMock.On("GetTransactionByUserID", mock.Anything, expectedUser.ID).Return(nil, expectedError)
	repoMock.On("RollbackTx", mock.Anything, txMock).Return(nil)

	loggerMock.On("Errorw",
		"Failed to get transactions",
		"userID", expectedUser.ID,
		"error", expectedError,
	).Return()

	service := NewUserService(repoMock, loggerMock)

	userInfo, err := service.GetInfo(context.Background(), expectedUser.ID)

	assert.Error(t, err)
	assert.Nil(t, userInfo)

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get transactions",
		"userID", expectedUser.ID,
		"error", expectedError,
	)

	repoMock.AssertCalled(t, "BeginTx", mock.Anything)
	repoMock.AssertCalled(t, "GetUserByID", mock.Anything, expectedUser.ID)
	repoMock.AssertCalled(t, "GetPurchaseByUserID", mock.Anything, expectedUser.ID)
	repoMock.AssertCalled(t, "GetTransactionByUserID", mock.Anything, expectedUser.ID)
	repoMock.AssertCalled(t, "RollbackTx", mock.Anything, txMock)
}

func TestGetInfo_Error_CommitTx(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)

	expectedCommitError := assert.AnError
	expectedUser := &models.User{
		ID:       1,
		Username: "user1",
		Balance:  100,
	}
	expectedPurchases := []*models.Purchase{
		{ID: 1, UserID: 1, MerchID: 10, CreatedAt: time.Now()},
	}
	expectedTransactions := []*models.Transaction{
		{ID: 1, SenderID: 1, ReceiverID: 2, Amount: 50, CreatedAt: time.Now()},
	}

	repoMock.On("BeginTx", mock.Anything).Return(txMock, nil)
	repoMock.On("GetUserByID", mock.Anything, expectedUser.ID).Return(expectedUser, nil)
	repoMock.On("GetPurchaseByUserID", mock.Anything, expectedUser.ID).Return(expectedPurchases, nil)
	repoMock.On("GetTransactionByUserID", mock.Anything, expectedUser.ID).Return(expectedTransactions, nil)
	repoMock.On("CommitTx", mock.Anything, txMock).Return(expectedCommitError)
	repoMock.On("RollbackTx", mock.Anything, txMock).Return(nil)

	loggerMock.On("Errorw",
		"Failed to commit transaction",
		"userID", expectedUser.ID,
		"error", expectedCommitError,
	).Return()

	service := NewUserService(repoMock, loggerMock)
	userInfo, err := service.GetInfo(context.Background(), expectedUser.ID)

	assert.Error(t, err)
	assert.Nil(t, userInfo)

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to commit transaction",
		"userID", expectedUser.ID,
		"error", expectedCommitError,
	)

	repoMock.AssertCalled(t, "BeginTx", mock.Anything)
	repoMock.AssertCalled(t, "CommitTx", mock.Anything, txMock)
	repoMock.AssertCalled(t, "RollbackTx", mock.Anything, txMock)
}

func TestGetInfo_RollbackTxError(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txMock := mockRepo.NewTxManager(t)

	userID := 1
	expectedUserError := errors.New("get user error")
	rollbackErr := errors.New("rollback error")

	repoMock.On("BeginTx", mock.Anything).Return(txMock, nil)
	repoMock.On("GetUserByID", mock.Anything, userID).Return(nil, expectedUserError)
	repoMock.On("RollbackTx", mock.Anything, txMock).Return(rollbackErr)

	loggerMock.On("Errorw",
		"Failed to get user info",
		"userID", userID,
		"error", expectedUserError,
	).Return()
	loggerMock.On("Errorw",
		"Failed to rollback transaction",
		"userID", userID,
		"error", rollbackErr,
	).Return()

	service := NewUserService(repoMock, loggerMock)
	userInfo, err := service.GetInfo(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, userInfo)
	assert.Contains(t, err.Error(), "get user error")

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get user info",
		"userID", userID,
		"error", expectedUserError,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Failed to rollback transaction",
		"userID", userID,
		"error", rollbackErr,
	)
	repoMock.AssertCalled(t, "BeginTx", mock.Anything)
	repoMock.AssertCalled(t, "GetUserByID", mock.Anything, userID)
	repoMock.AssertCalled(t, "RollbackTx", mock.Anything, txMock)
}
