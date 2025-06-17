package util

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken[T jwt.Claims](claims T) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken[T jwt.Claims](tokenStr string) (T, error) {
	var claims T
	token, err := jwt.ParseWithClaims(tokenStr, claims,
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

	if err != nil {
		var zero T
		return zero, err
	}

	if !token.Valid {
		var zero T
		return zero, errors.New("invalid token")
	}
	return claims, nil
}
