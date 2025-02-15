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

func TestBuyMerch_InvalidMerchID(t *testing.T) {
	mockRepo := new(mocks.Repository)
	service := NewMerchService(mockRepo, mockRepo, mockRepo)

	userID := uuid.New()
	merchID := 9999 // Несуществующий ID

	mockRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID, Balance: 100}, nil)
	mockRepo.On("GetMerchByID", mock.Anything, merchID).Return(nil, errors.New("товар не найден"))

	err := service.BuyMerch(context.TODO(), userID, merchID)

	assert.Error(t, err)
	assert.Equal(t, "товар не найден", err.Error())
}

func TestGetUserPurchases_EmptyList(t *testing.T) {
	mockRepo := new(mocks.Repository)
	service := NewMerchService(mockRepo, mockRepo, mockRepo)

	userID := uuid.New()

	mockRepo.On("GetPurchasesByUser", mock.Anything, userID).Return([]models.Purchase{}, nil)

	purchases, err := service.GetUserPurchases(context.TODO(), userID)

	assert.NoError(t, err)
	assert.Empty(t, purchases)
}

func TestGetMerchByName_Found(t *testing.T) {
	mockRepo := new(mocks.Repository)
	service := NewMerchService(mockRepo, mockRepo, mockRepo)

	merchName := "T-Shirt"
	expectedMerch := &models.Merch{ID: 1, Name: merchName, Price: 500}

	mockRepo.On("GetMerchByName", mock.Anything, merchName).Return(expectedMerch, nil)

	merch, err := service.GetMerchByName(context.TODO(), merchName)

	assert.NoError(t, err)
	assert.NotNil(t, merch)
	assert.Equal(t, merchName, merch.Name)
}
