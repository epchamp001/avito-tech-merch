package integration

import (
	"avito-tech-merch/internal/models"
	"avito-tech-merch/internal/storage/mocks"
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestBuyMerch(t *testing.T) {
	userID := uuid.New()
	merchID := 1
	merchPrice := 500

	mockRepo := new(mocks.Repository)

	mockRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID, Balance: 1000}, nil)
	mockRepo.On("GetMerchByID", mock.Anything, merchID).Return(&models.Merch{ID: merchID, Name: "Test Merch", Price: merchPrice}, nil)
	mockRepo.On("UpdateBalance", mock.Anything, userID, 500).Return(nil) // 1000 - 500 = 500
	mockRepo.On("CreatePurchase", mock.Anything, mock.AnythingOfType("*models.Purchase")).Return(nil)

	payload := map[string]interface{}{
		"user_id":  userID.String(),
		"merch_id": merchID,
	}

	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "http://localhost:8080/api/buy/test-item", bytes.NewReader(payloadBytes))
	assert.NoError(t, err)

	token := GetJWTToken(t)
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}
