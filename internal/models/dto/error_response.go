package dto

// ErrorResponse400 Response with an error
// @Description The standard API error format for 400 Bad Request
type ErrorResponse400 struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"invalid request"`
}

// ErrorResponse500 Response with an error
// @Description The standard API error format for 500 Internal Server Error
type ErrorResponse500 struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"server error"`
}

// ErrorResponseInvalidCredentials401 Response with an error
// @Description The standard API error format for 401 Unauthorized (invalid credentials)
type ErrorResponseInvalidCredentials401 struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"message" example:"invalid credentials"`
}

// ErrorResponseUnauthorized401 Response with an error
// @Description The standard API error format for 401 Unauthorized (general unauthorized error)
type ErrorResponseUnauthorized401 struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"message" example:"unauthorized"`
}

// ErrorJWTMissingToken represents an error response when the JWT token is missing
// @Description The standard API error format for 401 Unauthorized (missing token)
type ErrorJWTMissingToken struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"message" example:"missing token"`
}

// ErrorJWTInvalidToken represents an error response when the JWT token is invalid
// @Description The standard API error format for 401 Unauthorized (invalid token)
type ErrorJWTInvalidToken struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"message" example:"invalid token"`
}
