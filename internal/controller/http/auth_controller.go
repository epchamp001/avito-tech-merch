package http

import (
	"avito-tech-merch/internal/models/dto"
	"avito-tech-merch/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authController struct {
	service service.Service
}

func NewAuthController(service service.Service) AuthController {
	return &authController{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username and password, returns a JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 200 {object} dto.AuthResponse "JWT token"
// @Failure 400 {object} dto.ErrorResponse400 "Invalid request"
// @Failure 500 {object} dto.ErrorResponse500 "Internal server error"
// @Router /auth/register [post]
func (c *authController) Register(ctx *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse400{Code: 400, Message: "invalid request"})
		return
	}

	token, err := c.service.Register(ctx, request.Username, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse500{Code: 500, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.AuthResponse{Token: token})
}

// Login godoc
// @Summary Login a user
// @Description Login a user with username and password, returns a JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body dto.LoginRequest true "User login data"
// @Success 200 {object} dto.AuthResponse "JWT token"
// @Failure 400 {object} dto.ErrorResponse400 "Invalid request"
// @Failure 401 {object} dto.ErrorResponseInvalidCredentials401 "Invalid credentials"
// @Router /auth/login [post]
func (c *authController) Login(ctx *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse400{Code: 400, Message: "invalid request"})
		return
	}

	token, err := c.service.Login(ctx, request.Username, request.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponseInvalidCredentials401{Code: 401, Message: "invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, dto.AuthResponse{Token: token})
}
