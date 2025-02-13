package http

import (
	"avito-tech-merch/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type userController struct {
	service service.Service
}

func NewUserController(service service.Service) UserController {
	return &userController{service: service}
}

func (u *userController) GetUserInfo(c *gin.Context) {
	userID := c.GetString("userID")
	id, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
		return
	}

	info, err := u.service.GetUserInfo(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных"})
		return
	}

	c.JSON(http.StatusOK, info)
}

func (u *userController) TransferCoins(c *gin.Context) {
	var req struct {
		ToUser string `json:"toUser" binding:"required"`
		Amount int    `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	senderID := c.GetString("userID")
	senderUUID, err := uuid.Parse(senderID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
		return
	}

	receiver, err := u.service.GetUserByUsername(c.Request.Context(), req.ToUser)
	if err != nil || receiver == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Получатель не найден"})
		return
	}

	err = u.service.TransferCoins(c.Request.Context(), senderUUID, receiver.ID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Перевод выполнен"})
}
