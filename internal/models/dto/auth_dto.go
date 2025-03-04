package dto

// AuthRequest The structure of user credentials
// @Description Data for login or registration
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

// LoginRequest Login request
// @Description Data for user login
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"epchamp001"`
	Password string `json:"password" binding:"required" example:"strongpassword123"`
}

// AuthResponse Response for successful authentication or registration
// @Description Response containing JWT token
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2Njg5NTEwMTcsInN1YiI6ImpvaG5AZG9lLmNvbSJ9.Q3k6yMFYtuzPyjoZYpIHibJQPey29QWmlHfwS2A3keM"`
}
