package model

import "github.com/golang-jwt/jwt/v5"

type CustomerToken struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type SellerToken struct {
}
