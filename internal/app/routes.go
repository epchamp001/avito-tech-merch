package app

import (
	"avito-tech-merch/internal/controller/http"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, auth http.AuthController, user http.UserController, merch http.MerchController, secretKey string) {
	api := router.Group("/api")
	{
		api.POST("/auth", auth.AuthHandler)

		protected := api.Group("/")
		protected.Use(http.AuthMiddleware(secretKey))
		{
			protected.GET("/info", user.GetUserInfo)
			protected.POST("/sendCoin", user.TransferCoins)
			protected.POST("/buy/:item", merch.BuyMerch)
		}
	}
}
