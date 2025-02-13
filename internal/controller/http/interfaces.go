package http

import (
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	AuthHandler(c *gin.Context)
}

type UserController interface {
	GetUserInfo(c *gin.Context)
	TransferCoins(c *gin.Context)
}

type MerchController interface {
	BuyMerch(c *gin.Context)
}
