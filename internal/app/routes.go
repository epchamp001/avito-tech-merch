package app

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/info", func(c *gin.Context) {})
	}
}
