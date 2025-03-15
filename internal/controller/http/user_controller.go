package http

import (
	"avito-tech-merch/internal/models/dto"
	"avito-tech-merch/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type userController struct {
	service service.Service
}

func NewUserController(service service.Service) UserController {
	return &userController{service: service}
}

// GetInfo godoc
// @Summary Get user information
// @Security BearerAuth
// @Description Fetches user information based on the userID from the context
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.UserInfoResponse "User information"
// @Failure 401 {object} dto.ErrorResponseUnauthorized401 "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse500 "Internal server error"
// @Router /info [get]
func (c *userController) GetInfo(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponseUnauthorized401{Code: 401, Message: "unauthorized"})
		return
	}

	info, err := c.service.GetInfo(ctx, userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse500{Code: 500, Message: err.Error()})
		return
	}

	infoDTO := dto.MapUserInfoResponseToDTO(info)

	ctx.JSON(http.StatusOK, infoDTO)
}
