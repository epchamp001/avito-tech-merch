package http

import (
	"avito-tech-merch/internal/models/dto"
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

// BuyMerch godoc
// @Summary Purchase a merchandise item
// @Description Allows a user to purchase a merchandise item by specifying the item ID in the URL
// @Tags purchase
// @Accept  json
// @Produce  json
// @Param item path string true "Item ID to purchase" example:"cup"
// @Success 200 {object} dto.PurchaseSuccessResponse "Purchase successful"
// @Failure 400 {object} dto.ErrorResponse400 "Bad request (item is required)"
// @Failure 401 {object} dto.ErrorResponseUnauthorized401 "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse500 "Internal server error"
// @Router /merch/buy/:item [post]
func (c *purchaseController) BuyMerch(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponseUnauthorized401{Code: 401, Message: "unauthorized"})
		return
	}

	item := ctx.Param("item")
	if item == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse400{Code: 400, Message: "item is required"})
		return
	}

	err := c.service.PurchaseMerch(ctx, userID.(int), item)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse500{Code: 500, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.PurchaseSuccessResponse{Message: "purchase successful"})
}
