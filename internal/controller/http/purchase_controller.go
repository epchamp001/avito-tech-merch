package http

import (
	"avito-tech-merch/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type purchaseController struct {
	service service.Service
}

func NewPurchaseController(service service.Service) PurchaseController {
	return &purchaseController{service: service}
}

func (c *purchaseController) BuyMerch(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	item := ctx.Param("item")
	if item == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "item is required"})
		return
	}

	err := c.service.PurchaseMerch(ctx, userID.(int), item)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "purchase successful"})
}
