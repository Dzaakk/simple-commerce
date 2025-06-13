package util

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(claims model.CustomerToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenStr string) (*model.CustomerToken, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.CustomerToken{},
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.CustomerToken)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
