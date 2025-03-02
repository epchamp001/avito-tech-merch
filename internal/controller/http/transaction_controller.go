package http

import (
	"avito-tech-merch/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transactionController struct {
	service service.Service
}

func NewTransactionController(service service.Service) TransactionController {
	return &transactionController{service: service}
}

func (c *transactionController) SendCoin(ctx *gin.Context) {
	senderID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request struct {
		ReceiverID int `json:"receiver_id" binding:"required"`
		Amount     int `json:"amount" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := c.service.TransferCoins(ctx, senderID.(int), request.ReceiverID, request.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "coins transferred successfully"})
}
