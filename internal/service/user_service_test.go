package service

import (
	"avito-tech-merch/internal/models"
	mockRepo "avito-tech-merch/internal/storage/db/mock"
	"avito-tech-merch/internal/storage/db/postgres"
	mockLog "avito-tech-merch/pkg/logger/mock"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestGetInfo_Success(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)

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

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelReadCommitted, postgres.AccessModeReadOnly, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			// Вызываем функцию транзакционного блока
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(nil)

	repoMock.On("GetUserByID", mock.Anything, expectedUser.ID).Return(expectedUser, nil)
	repoMock.On("GetPurchaseByUserID", mock.Anything, expectedUser.ID).Return(expectedPurchases, nil)
	repoMock.On("GetTransactionByUserID", mock.Anything, expectedUser.ID).Return(expectedTransactions, nil)

	service := NewUserService(repoMock, loggerMock, txManagerMock)

	userInfo, err := service.GetInfo(context.Background(), expectedUser.ID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, userInfo.UserID)
	assert.Equal(t, expectedUser.Username, userInfo.Username)
	assert.Equal(t, expectedUser.Balance, userInfo.Balance)
	assert.Equal(t, expectedPurchases, userInfo.Purchases)
	assert.Equal(t, expectedTransactions, userInfo.Transactions)

	repoMock.AssertExpectations(t)
	txManagerMock.AssertExpectations(t)
}

func TestGetInfo_Error_GetUser(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)

	expectedError := assert.AnError
	userID := 1

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelReadCommitted, postgres.AccessModeReadOnly, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(expectedError)

	repoMock.On("GetUserByID", mock.Anything, userID).Return(nil, expectedError)

	loggerMock.
		On("Errorw",
			"Failed to get user info",
			"userID", userID,
			"error", expectedError,
		).Return()
	loggerMock.
		On("Errorw",
			"Error getting user info",
			"error", expectedError,
		).Return()

	service := NewUserService(repoMock, loggerMock, txManagerMock)

	userInfo, err := service.GetInfo(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, userInfo)

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get user info",
		"userID", userID,
		"error", expectedError,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Error getting user info",
		"error", expectedError,
	)

	repoMock.AssertCalled(t, "GetUserByID", mock.Anything, userID)
	txManagerMock.AssertExpectations(t)
}

func TestGetInfo_Error_GetPurchases(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)

	expectedError := assert.AnError
	expectedUser := &models.User{
		ID:       1,
		Username: "user1",
		Balance:  100,
	}

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelReadCommitted, postgres.AccessModeReadOnly, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(expectedError)

	repoMock.On("GetUserByID", mock.Anything, expectedUser.ID).Return(expectedUser, nil)
	repoMock.On("GetPurchaseByUserID", mock.Anything, expectedUser.ID).Return(nil, expectedError)

	loggerMock.
		On("Errorw",
			"Failed to get purchases",
			"userID", expectedUser.ID,
			"error", expectedError,
		).Return()
	loggerMock.
		On("Errorw",
			"Error getting user info",
			"error", expectedError,
		).Return()

	service := NewUserService(repoMock, loggerMock, txManagerMock)

	userInfo, err := service.GetInfo(context.Background(), expectedUser.ID)

	assert.Error(t, err)
	assert.Nil(t, userInfo)

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get purchases",
		"userID", expectedUser.ID,
		"error", expectedError,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Error getting user info",
		"error", expectedError,
	)

	repoMock.AssertCalled(t, "GetUserByID", mock.Anything, expectedUser.ID)
	repoMock.AssertCalled(t, "GetPurchaseByUserID", mock.Anything, expectedUser.ID)
	txManagerMock.AssertExpectations(t)
}

func TestGetInfo_Error_GetTransactions(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)

	expectedError := assert.AnError
	expectedUser := &models.User{
		ID:       1,
		Username: "user1",
		Balance:  100,
	}
	expectedPurchases := []*models.Purchase{
		{ID: 1, UserID: 1, MerchID: 10, CreatedAt: time.Now()},
	}

	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelReadCommitted, postgres.AccessModeReadOnly, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(expectedError)

	repoMock.On("GetUserByID", mock.Anything, expectedUser.ID).Return(expectedUser, nil)
	repoMock.On("GetPurchaseByUserID", mock.Anything, expectedUser.ID).Return(expectedPurchases, nil)
	repoMock.On("GetTransactionByUserID", mock.Anything, expectedUser.ID).Return(nil, expectedError)

	loggerMock.
		On("Errorw",
			"Failed to get transactions",
			"userID", expectedUser.ID,
			"error", expectedError,
		).Return()
	loggerMock.
		On("Errorw",
			"Error getting user info",
			"error", expectedError,
		).Return()

	service := NewUserService(repoMock, loggerMock, txManagerMock)

	userInfo, err := service.GetInfo(context.Background(), expectedUser.ID)

	assert.Error(t, err)
	assert.Nil(t, userInfo)

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to get transactions",
		"userID", expectedUser.ID,
		"error", expectedError,
	)
	loggerMock.AssertCalled(t, "Errorw",
		"Error getting user info",
		"error", expectedError,
	)

	repoMock.AssertCalled(t, "GetUserByID", mock.Anything, expectedUser.ID)
	repoMock.AssertCalled(t, "GetPurchaseByUserID", mock.Anything, expectedUser.ID)
	repoMock.AssertCalled(t, "GetTransactionByUserID", mock.Anything, expectedUser.ID)
	txManagerMock.AssertExpectations(t)
}

func TestGetInfo_Error_FromTxManager(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)
	txManagerMock := mockRepo.NewTxManager(t)

	expectedError := assert.AnError
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

	// Даже если репозиторий отдает корректные данные, txManager возвращает ошибку (например, на commit)
	txManagerMock.
		On("WithTx", mock.Anything, postgres.IsolationLevelReadCommitted, postgres.AccessModeReadOnly, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(3).(func(context.Context) error)
			_ = fn(context.Background())
		}).
		Return(expectedError)

	repoMock.On("GetUserByID", mock.Anything, expectedUser.ID).Return(expectedUser, nil)
	repoMock.On("GetPurchaseByUserID", mock.Anything, expectedUser.ID).Return(expectedPurchases, nil)
	repoMock.On("GetTransactionByUserID", mock.Anything, expectedUser.ID).Return(expectedTransactions, nil)

	loggerMock.
		On("Errorw",
			"Error getting user info",
			"error", expectedError,
		).Return()

	service := NewUserService(repoMock, loggerMock, txManagerMock)
	userInfo, err := service.GetInfo(context.Background(), expectedUser.ID)

	assert.Error(t, err)
	assert.Nil(t, userInfo)

	loggerMock.AssertCalled(t, "Errorw",
		"Error getting user info",
		"error", expectedError,
	)

	repoMock.AssertCalled(t, "GetUserByID", mock.Anything, expectedUser.ID)
	repoMock.AssertCalled(t, "GetPurchaseByUserID", mock.Anything, expectedUser.ID)
	repoMock.AssertCalled(t, "GetTransactionByUserID", mock.Anything, expectedUser.ID)
	txManagerMock.AssertExpectations(t)
}
