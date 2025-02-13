package http

import (
	"avito-tech-merch/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type merchController struct {
	service service.MerchService
}

func NewMerchController(service service.MerchService) MerchController {
	return &merchController{service: service}
}

func (m *merchController) BuyMerch(c *gin.Context) {
	itemName := c.Param("item")

	userID := c.GetString("userID")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
		return
	}

	merch, err := m.service.GetMerchByName(c.Request.Context(), itemName)
	if err != nil || merch == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
		return
	}

	err = m.service.BuyMerch(c.Request.Context(), userUUID, merch.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Покупка успешна"})
}
