package service

import (
	"avito-tech-merch/internal/models"
	mockRepo "avito-tech-merch/internal/storage/db/mock"
	mockLog "avito-tech-merch/pkg/logger/mock"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListMerch(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)

	expectedMerch := []*models.Merch{
		{ID: 1, Name: "Merch 1"},
		{ID: 2, Name: "Merch 2"},
	}

	ctx := context.Background()

	repoMock.On("GetAllMerch", ctx).Return(expectedMerch, nil)

	service := NewMerchService(repoMock, loggerMock)

	merchList, err := service.ListMerch(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedMerch, merchList)
}

func TestListMerch_Error(t *testing.T) {
	repoMock := mockRepo.NewRepository(t)
	loggerMock := mockLog.NewLogger(t)

	ctx := context.Background()

	expectedError := assert.AnError

	repoMock.On("GetAllMerch", ctx).Return(nil, expectedError)

	loggerMock.On("Errorw",
		"Failed to fetch merch list",
		"error", expectedError,
	).Return()

	service := NewMerchService(repoMock, loggerMock)

	merchList, err := service.ListMerch(ctx)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	assert.Nil(t, merchList)

	repoMock.AssertCalled(t, "GetAllMerch", ctx)

	loggerMock.AssertCalled(t, "Errorw",
		"Failed to fetch merch list",
		"error", expectedError,
	)
}

func TestGetMerch(t *testing.T) {
	tests := []struct {
		name          string
		merchID       int
		expectedMerch *models.Merch
		mockError     error
		expectedError error
	}{
		{
			name:    "Successfully fetch merch",
			merchID: 1,
			expectedMerch: &models.Merch{
				ID:   1,
				Name: "Merch 1",
			},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name:          "Error fetching merch - not found",
			merchID:       999,
			expectedMerch: nil,
			mockError:     errors.New("merch not found"),
			expectedError: errors.New("merch not found"),
		},
		{
			name:          "Error fetching merch - general error",
			merchID:       1,
			expectedMerch: nil,
			mockError:     assert.AnError,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mockRepo.NewRepository(t)
			loggerMock := mockLog.NewLogger(t)

			ctx := context.Background()

			repoMock.On("GetMerchByID", ctx, tt.merchID).Return(tt.expectedMerch, tt.mockError)

			if tt.mockError != nil {
				loggerMock.On("Errorw", "Failed to fetch merch", "merchID", tt.merchID, "error", tt.mockError).Return()
			}

			service := NewMerchService(repoMock, loggerMock)

			merch, err := service.GetMerch(ctx, tt.merchID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMerch, merch)
			}

			repoMock.AssertCalled(t, "GetMerchByID", ctx, tt.merchID)

			if tt.mockError != nil {
				loggerMock.AssertCalled(t, "Errorw", "Failed to fetch merch", "merchID", tt.merchID, "error", tt.mockError)
			}
		})
	}
}
