package http

import (
	"avito-tech-merch/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен отсутствует"})
			c.Abort()
			return
		}

		userID, err := utils.ParseJWT(tokenString, secretKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
