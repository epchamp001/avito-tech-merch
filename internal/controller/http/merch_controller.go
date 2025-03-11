package http

import (
	"avito-tech-merch/internal/models/dto"
	"avito-tech-merch/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type merchController struct {
	service service.Service
}

func NewMerchController(service service.Service) MerchController {
	return &merchController{service: service}
}

// ListMerch godoc
// @Summary Get list of merchandise items
// @Description Fetches all merch items from the database and returns a list of merch
// @Tags merch
// @Accept  json
// @Produce  json
// @Success 200 {array} dto.MerchDTO "List of merchandise items"
// @Failure 500 {object} dto.ErrorResponse500 "Internal server error"
// @Router /merch [get]
func (c *merchController) ListMerch(ctx *gin.Context) {
	merchList, err := c.service.ListMerch(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse500{Code: 500, Message: err.Error()})
		return
	}

	merchListDTO := make([]*dto.MerchDTO, len(merchList))
	for i, merch := range merchList {
		merchListDTO[i] = dto.MapMerchToDTO(merch)
	}

	ctx.JSON(http.StatusOK, merchListDTO)
}
