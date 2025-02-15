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

func TestTransferCoins(t *testing.T) {
	senderID := uuid.New()
	receiverID := uuid.New()
	amount := 200

	mockRepo := new(mocks.Repository)

	mockRepo.On("GetUserByID", mock.Anything, senderID).Return(&models.User{ID: senderID, Balance: 1000}, nil)
	mockRepo.On("GetUserByID", mock.Anything, receiverID).Return(&models.User{ID: receiverID, Balance: 500}, nil)
	mockRepo.On("UpdateBalance", mock.Anything, senderID, 800).Return(nil)   // 1000 - 200 = 800
	mockRepo.On("UpdateBalance", mock.Anything, receiverID, 700).Return(nil) // 500 + 200 = 700
	mockRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*models.Transaction")).Return(nil)

	payload := map[string]interface{}{
		"sender_id":   senderID.String(),
		"receiver_id": receiverID.String(),
		"amount":      amount,
	}

	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "http://localhost:8080/api/sendCoin", bytes.NewReader(payloadBytes))
	assert.NoError(t, err)

	token := GetJWTToken(t)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}
