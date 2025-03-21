package middleware

import (
	"avito-tech-merch/internal/models/dto"
	"avito-tech-merch/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JWTAuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorJWTMissingToken{
				Code:    401,
				Message: "missing token",
			})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		userID, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorJWTInvalidToken{
				Code:    401,
				Message: "invalid token",
			})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
