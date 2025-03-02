package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateToken(userID int, secretKey string, expirationTime int) (string, error) {
	now := time.Now()

	expiration := now.Add(time.Duration(expirationTime) * time.Second)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiration.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}
