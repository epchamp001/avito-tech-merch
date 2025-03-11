package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenService interface {
	GenerateToken(userID int, secretKey string, expirationTime int) (string, error)
	ParseJWTToken(tokenString, secretKey string) (int, error)
}

type TokenServiceImpl struct {
}

func NewTokenService() TokenService {
	return &TokenServiceImpl{}
}

func (t *TokenServiceImpl) GenerateToken(userID int, secretKey string, expirationTime int) (string, error) {
	now := time.Now()

	expiration := now.Add(time.Duration(expirationTime) * time.Second)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiration.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

func (t *TokenServiceImpl) ParseJWTToken(tokenString, secretKey string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, fmt.Errorf("invalid token")
}
