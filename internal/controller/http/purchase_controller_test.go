package http

import (
	"avito-tech-merch/internal/models/dto"
	mockServ "avito-tech-merch/internal/service/mock"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPurchaseController_BuyMerch_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)
	controller := NewPurchaseController(mockService)
	router := gin.New()
	router.POST("/merch/buy/:item", controller.BuyMerch)

	// userID не установлен в контексте
	req, _ := http.NewRequest("POST", "/merch/buy/cup", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	var resp dto.ErrorResponseUnauthorized401
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.Code)
	assert.Equal(t, "unauthorized", resp.Message)
}

func TestPurchaseController_BuyMerch_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)
	controller := NewPurchaseController(mockService)
	router := gin.New()

	// Регистрируем middleware для установки userID до регистрации маршрута
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.POST("/merch/buy/:item", controller.BuyMerch)

	// Ожидаем, что сервис вернет ошибку при покупке
	serviceErr := errors.New("purchase failed")
	mockService.On("PurchaseMerch", mock.Anything, 1, "cup").Return(serviceErr).Once()

	req, _ := http.NewRequest("POST", "/merch/buy/cup", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var resp dto.ErrorResponse500
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.Code)
	assert.Equal(t, serviceErr.Error(), resp.Message)

	mockService.AssertCalled(t, "PurchaseMerch", mock.Anything, 1, "cup")
}

func TestPurchaseController_BuyMerch_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)
	controller := NewPurchaseController(mockService)
	router := gin.New()

	// Регистрируем middleware для установки userID до объявления маршрута
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.POST("/merch/buy/:item", controller.BuyMerch)

	mockService.On("PurchaseMerch", mock.Anything, 1, "cup").Return(nil).Once()

	req, _ := http.NewRequest("POST", "/merch/buy/cup", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp dto.PurchaseSuccessResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "purchase successful", resp.Message)

	mockService.AssertCalled(t, "PurchaseMerch", mock.Anything, 1, "cup")
}
