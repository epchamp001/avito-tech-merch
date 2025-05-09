package app

import (
	controller "avito-tech-merch/internal/controller/http"
	"avito-tech-merch/internal/controller/http/middleware"
	"avito-tech-merch/internal/metrics"
	"avito-tech-merch/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, controller controller.Controller, authService service.Service) {
	authMiddleware := middleware.JWTAuthMiddleware(authService)

	router.Use(metrics.GinPrometheusMiddleware())

	api := router.Group("/api")
	{
		api.POST("/auth/register", controller.Register)

		api.POST("/auth/login", controller.Login)

		protected := api.Group("/")
		protected.Use(authMiddleware)
		{
			protected.GET("/info", controller.GetInfo)
			protected.POST("/send-coin", controller.SendCoin)
			protected.GET("/merch", controller.ListMerch)
			protected.POST("/merch/buy/:item", controller.BuyMerch)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
