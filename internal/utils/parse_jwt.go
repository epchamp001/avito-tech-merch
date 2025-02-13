package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

func ParseJWT(tokenString, secretKey string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("неверный формат токена")
	}

	return userID, nil
}
