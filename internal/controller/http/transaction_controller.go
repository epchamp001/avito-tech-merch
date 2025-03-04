package http

import (
	"avito-tech-merch/internal/models/dto"
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

// SendCoin godoc
// @Summary Transfer coins between users
// @Description Allows a user to send coins to another user by specifying the receiver ID and the amount
// @Tags transaction
// @Accept  json
// @Produce  json
// @Param receiver_id body int true "Receiver user ID" example:"2"
// @Param amount body int true "Amount of coins to transfer" example:"100"
// @Success 200 {object} dto.TransferSuccessResponse "Coins transferred successfully"
// @Failure 400 {object} dto.ErrorResponse400 "Invalid request (missing or invalid data)"
// @Failure 401 {object} dto.ErrorResponseUnauthorized401 "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse500 "Internal server error"
// @Router /send-coin [post]
func (c *transactionController) SendCoin(ctx *gin.Context) {
	senderID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponseUnauthorized401{Code: 401, Message: "unauthorized"})
		return
	}

	var request struct {
		ReceiverID int `json:"receiver_id" binding:"required"`
		Amount     int `json:"amount" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse400{Code: 400, Message: "invalid request"})
		return
	}

	err := c.service.TransferCoins(ctx, senderID.(int), request.ReceiverID, request.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse500{Code: 500, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.TransferSuccessResponse{Message: "coins transferred successfully"})
}
