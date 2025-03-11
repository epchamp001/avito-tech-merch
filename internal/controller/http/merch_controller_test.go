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

func TestMerchController_ListMerch_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)

	controller := NewMerchController(mockService)
	router := gin.New()
	router.GET("/merch", controller.ListMerch)

	expectedErr := errors.New("service error")
	mockService.On("ListMerch", mock.Anything).Return(nil, expectedErr).Once()

	req, _ := http.NewRequest("GET", "/merch", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var errResp dto.ErrorResponse500
	err := json.Unmarshal(rec.Body.Bytes(), &errResp)
	assert.NoError(t, err)
	assert.Equal(t, 500, errResp.Code)
	assert.Equal(t, expectedErr.Error(), errResp.Message)

	mockService.AssertCalled(t, "ListMerch", mock.Anything)
}

func TestMerchController_ListMerch_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := mockServ.NewService(t)
	controller := NewMerchController(mockService)
	router := gin.New()
	router.GET("/merch", controller.ListMerch)

	merchList := []*models.Merch{
		{ID: 1, Name: "cup", Price: 20},
		{ID: 2, Name: "t-shirt", Price: 50},
	}

	mockService.On("ListMerch", mock.Anything).Return(merchList, nil).Once()

	req, _ := http.NewRequest("GET", "/merch", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	expectedDTO := make([]*dto.MerchDTO, len(merchList))
	for i, m := range merchList {
		expectedDTO[i] = dto.MapMerchToDTO(m)
	}

	var resp []*dto.MerchDTO
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedDTO, resp)

	mockService.AssertCalled(t, "ListMerch", mock.Anything)
}
