package dto

import "Dzaakk/simple-commerce/package/constant"

type RegisterCustomerRequest struct {
	Email    string
	Password string
	FullName string
	Phone    string
}

type RegisterSellerRequest struct {
	Email    string
	Password string
	FullName string
	Phone    string
	ShopName string
}

type LoginRequest struct {
	Email    string
	Password string
	UserType constant.UserType
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

type RefreshTokenResponse struct {
	AccessToken string
	ExpiresIn   int
}
