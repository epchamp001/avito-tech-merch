package service

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/mocks"
	"avito-tech-merch/internal/utils"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRegisterUser_AlreadyExists(t *testing.T) {
	mockRepo := new(mocks.Repository)
	service := NewUserService(mockRepo)

	username := "testuser"
	password := "securepassword"

	// Юзер уже существует
	existingUser := &models.User{ID: uuid.New(), Username: username}
	mockRepo.On("GetUserByUsername", mock.Anything, username).Return(existingUser, nil)

	user, err := service.RegisterUser(context.TODO(), username, password)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "пользователь уже существует", err.Error())
}

func TestAuthenticateUser_SuccessfulLogin(t *testing.T) {
	mockRepo := new(mocks.Repository)
	service := NewUserService(mockRepo)

	username := "testuser"
	password := "password"

	// Генерируем корректный хеш пароля
	hashedPassword, _ := utils.HashPassword(password)

	user := &models.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: hashedPassword,
	}

	mockRepo.On("GetUserByUsername", mock.Anything, username).Return(user, nil)

	authUser, err := service.AuthenticateUser(context.TODO(), username, password)

	assert.NoError(t, err)
	assert.NotNil(t, authUser)
	assert.Equal(t, username, authUser.Username)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_UserExists(t *testing.T) {
	mockRepo := new(mocks.Repository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	expectedUser := &models.User{ID: userID, Username: "testuser"}

	mockRepo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)

	user, err := service.GetUserByID(context.TODO(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.Username, user.Username)
}

func TestGetUserByUsername_Found(t *testing.T) {
	mockRepo := new(mocks.Repository)
	service := NewUserService(mockRepo)

	username := "testuser"
	expectedUser := &models.User{ID: uuid.New(), Username: username}

	mockRepo.On("GetUserByUsername", mock.Anything, username).Return(expectedUser, nil)

	user, err := service.GetUserByUsername(context.TODO(), username)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
}

func TestUpdateBalance_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	mockRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID, Balance: 100}, nil)
	mockRepo.On("UpdateBalance", mock.Anything, userID, 200).Return(nil)

	err := service.UpdateBalance(context.TODO(), userID, 100)

	assert.NoError(t, err)
}
