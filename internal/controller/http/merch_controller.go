package http

import (
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

func (c *merchController) ListMerch(ctx *gin.Context) {
	merchList, err := c.service.ListMerch(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, merchList)
}
