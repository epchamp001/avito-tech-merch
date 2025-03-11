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

func TestAuthController_Register_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)

	controller := NewAuthController(mockService)
	router := gin.New()
	router.POST("/auth/register", controller.Register)

	// Отправляем некорректный JSON (отсутствует поле password)
	reqBody := `{"username": "epchamp001"}`
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer([]byte(reqBody)))
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

func TestAuthController_Register_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)

	controller := NewAuthController(mockService)
	router := gin.New()
	router.POST("/auth/register", controller.Register)

	reqData := dto.RegisterRequest{
		Username: "epchamp001",
		Password: "strongpassword123",
	}
	reqBody, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	serviceErr := errors.New("registration failed")
	mockService.On("Register", mock.Anything, "epchamp001", "strongpassword123").Return("", serviceErr).Once()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var resp dto.ErrorResponse500
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.Code)
	assert.Equal(t, serviceErr.Error(), resp.Message)

	mockService.AssertCalled(t, "Register", mock.Anything, "epchamp001", "strongpassword123")
}

func TestAuthController_Register_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)

	controller := NewAuthController(mockService)
	router := gin.New()
	router.POST("/auth/register", controller.Register)

	reqData := dto.RegisterRequest{
		Username: "epchamp001",
		Password: "strongpassword123",
	}
	reqBody, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	expectedToken := "jwt-token"
	mockService.On("Register", mock.Anything, "epchamp001", "strongpassword123").Return(expectedToken, nil).Once()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp dto.AuthResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, resp.Token)

	mockService.AssertCalled(t, "Register", mock.Anything, "epchamp001", "strongpassword123")
}

func TestAuthController_Login_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)

	controller := NewAuthController(mockService)
	router := gin.New()
	router.POST("/auth/login", controller.Login)

	// Отправляем некорректный JSON (отсутствует поле password)
	reqBody := `{"username": "epchamp001"}`
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte(reqBody)))
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

func TestAuthController_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)

	controller := NewAuthController(mockService)
	router := gin.New()
	router.POST("/auth/login", controller.Login)

	reqData := dto.LoginRequest{
		Username: "epchamp001",
		Password: "wrongpassword",
	}
	reqBody, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	mockService.On("Login", mock.Anything, "epchamp001", "wrongpassword").Return("", errors.New("invalid credentials")).Once()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	var resp dto.ErrorResponseInvalidCredentials401
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.Code)
	assert.Equal(t, "invalid credentials", resp.Message)

	mockService.AssertCalled(t, "Login", mock.Anything, "epchamp001", "wrongpassword")
}

func TestAuthController_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)

	controller := NewAuthController(mockService)
	router := gin.New()
	router.POST("/auth/login", controller.Login)

	reqData := dto.LoginRequest{
		Username: "epchamp001",
		Password: "strongpassword123",
	}
	reqBody, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	expectedToken := "jwt-token"
	mockService.On("Login", mock.Anything, "epchamp001", "strongpassword123").Return(expectedToken, nil).Once()

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp dto.AuthResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, resp.Token)

	mockService.AssertCalled(t, "Login", mock.Anything, "epchamp001", "strongpassword123")
}
