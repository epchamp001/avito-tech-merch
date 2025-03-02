package http

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	AuthController
	UserController
	MerchController
	PurchaseController
	TransactionController
}

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type UserController interface {
	GetInfo(ctx *gin.Context)
}

type MerchController interface {
	ListMerch(ctx *gin.Context)
}

type PurchaseController interface {
	BuyMerch(ctx *gin.Context)
}

type TransactionController interface {
	SendCoin(ctx *gin.Context)
}

type controller struct {
	AuthController
	UserController
	MerchController
	PurchaseController
	TransactionController
}

func NewController(
	auth AuthController,
	user UserController,
	merch MerchController,
	purchase PurchaseController,
	transaction TransactionController,
) Controller {
	return &controller{
		AuthController:        auth,
		UserController:        user,
		MerchController:       merch,
		PurchaseController:    purchase,
		TransactionController: transaction,
	}
}
