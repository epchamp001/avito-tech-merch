package http

import (
	"avito-tech-merch/internal/models"
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

func TestUserController_GetInfo_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)
	controller := NewUserController(mockService)
	router := gin.New()
	router.GET("/info", controller.GetInfo)

	// Не устанавливаем userID в контекст
	req, _ := http.NewRequest("GET", "/info", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	var resp dto.ErrorResponseUnauthorized401
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.Code)
	assert.Equal(t, "unauthorized", resp.Message)
}

func TestUserController_GetInfo_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)
	controller := NewUserController(mockService)
	router := gin.New()

	// Устанавливаем userID в контекст
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.GET("/info", controller.GetInfo)

	expectedErr := errors.New("service error")
	mockService.On("GetInfo", mock.Anything, 1).Return(nil, expectedErr).Once()

	req, _ := http.NewRequest("GET", "/info", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var resp dto.ErrorResponse500
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.Code)
	assert.Equal(t, expectedErr.Error(), resp.Message)

	mockService.AssertCalled(t, "GetInfo", mock.Anything, 1)
}

func TestUserController_GetInfo_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)
	controller := NewUserController(mockService)
	router := gin.New()

	// Устанавливаем userID в контекст
	router.Use(func(c *gin.Context) {
		c.Set("userID", 1)
		c.Next()
	})
	router.GET("/info", controller.GetInfo)

	// Создаем тестовый объект userInfo
	userInfo := &models.UserInfo{
		UserID:       1,
		Username:     "testuser",
		Balance:      1000,
		Purchases:    []*models.Purchase{},
		Transactions: []*models.Transaction{},
	}
	expectedDTO := dto.MapUserInfoResponseToDTO(userInfo)

	mockService.On("GetInfo", mock.Anything, 1).Return(userInfo, nil).Once()

	req, _ := http.NewRequest("GET", "/info", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp dto.UserInfoResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedDTO, &resp)

	mockService.AssertCalled(t, "GetInfo", mock.Anything, 1)
}
