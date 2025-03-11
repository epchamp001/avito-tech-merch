package http

import (
	"avito-tech-merch/internal/models/dto"
	mockServ "avito-tech-merch/internal/service/mock"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransactionController_SendCoin_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)
	controller := NewTransactionController(mockService)
	router := gin.New()
	router.POST("/send-coin", controller.SendCoin)

	// Не устанавливаем userID в контекст
	reqBody, _ := json.Marshal(dto.TransferRequest{
		ReceiverID: 2,
		Amount:     100,
	})
	req, _ := http.NewRequest("POST", "/send-coin", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	var resp dto.ErrorResponseUnauthorized401
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.Code)
	assert.Equal(t, "unauthorized", resp.Message)
}

func TestTransactionController_SendCoin_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)
	controller := NewTransactionController(mockService)
	router := gin.New()

	// Регистрируем middleware для установки userID
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.POST("/send-coin", controller.SendCoin)

	// Отправляем некорректный JSON (отсутствует поле amount)
	reqBody := `{"receiver_id": 2}`
	req, _ := http.NewRequest("POST", "/send-coin", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var resp dto.ErrorResponse400
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "invalid request", resp.Message)
}

func TestTransactionController_SendCoin_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)
	controller := NewTransactionController(mockService)
	router := gin.New()

	// Устанавливаем userID в контекст
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.POST("/send-coin", controller.SendCoin)

	// корректный запрос
	reqData := dto.TransferRequest{
		ReceiverID: 2,
		Amount:     100,
	}
	reqBody, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/send-coin", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	serviceErr := errors.New("transfer failed")
	// Ожидаем, что сервис вернет ошибку при вызове TransferCoins
	mockService.On("TransferCoins", mock.Anything, 1, 2, 100).Return(serviceErr).Once()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var resp dto.ErrorResponse500
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.Code)
	assert.Equal(t, serviceErr.Error(), resp.Message)

	mockService.AssertCalled(t, "TransferCoins", mock.Anything, 1, 2, 100)
}

func TestTransactionController_SendCoin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)
	controller := NewTransactionController(mockService)
	router := gin.New()

	// Устанавливаем userID в контекст
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.POST("/send-coin", controller.SendCoin)

	// Подготавливаем корректный запрос
	reqData := dto.TransferRequest{
		ReceiverID: 2,
		Amount:     100,
	}
	reqBody, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/send-coin", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Ожидаем успешный вызов сервиса
	mockService.On("TransferCoins", mock.Anything, 1, 2, 100).Return(nil).Once()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp dto.TransferSuccessResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "coins transferred successfully", resp.Message)

	mockService.AssertCalled(t, "TransferCoins", mock.Anything, 1, 2, 100)
}
