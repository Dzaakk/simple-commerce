package util

import (
	"Dzaakk/simple-commerce/internal/auth/model"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(secretKey []byte, claims model.CustomerToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
