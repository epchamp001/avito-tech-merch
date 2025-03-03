package dto

type AuthRequest struct {
	Username string `json:"username" binding:"required" example:"epchamp001"`
	Password string `json:"password" binding:"required" example:"strongpassword123"`
}

// RegisterRequest Registration request
// @Description Data for creating a new user
type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"epchamp001"`
	Password string `json:"password" binding:"required" example:"strongpassword123"`
}

// LoginRequest Authentication Request
// @Description Login information
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"epchamp001"`
	Password string `json:"password" binding:"required" example:"strongpassword123"`
}

// AuthResponse Successful response with a token
// @Description Contains a JWT token for authentication
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Njg5NTEwMTcsInN1YiI6ImpvaG5AZG9lLmNvbSJ9.Q3k6yMFYtuzPyjoZYpIHibJQPey29QWmlHfwS2A3keM"`
}

// ErrorResponse Response with an error
// @Description The standard API error format
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"invalid request"`
}
