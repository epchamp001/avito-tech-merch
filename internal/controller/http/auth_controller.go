package http

import (
	"avito-tech-merch/internal/service"
	"avito-tech-merch/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authController struct {
	service   service.Service
	secretKey string
}

func NewAuthController(service service.Service, secretKey string) AuthController {
	return &authController{service: service, secretKey: secretKey}
}

func (a *authController) AuthHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	user, err := a.service.AuthenticateUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateJWT(user.ID, a.secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
