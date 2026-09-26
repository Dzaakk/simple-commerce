package service

import (
	"context"

	"Dzaakk/simple-commerce/internal/auth/dto"
	"Dzaakk/simple-commerce/internal/auth/model"
	userdto "Dzaakk/simple-commerce/internal/user/dto"
	usermodel "Dzaakk/simple-commerce/internal/user/model"
)

type AuthService interface {
	RegisterCustomer(context.Context, *dto.RegisterCustomerRequest) error
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(context.Context, string) (*dto.RefreshTokenResponse, error)
	Logout(context.Context, string) error
}

type customerService interface {
	Create(context.Context, *userdto.RegisterCustomerRequest) (string, error)
	FindByEmail(context.Context, string) (*usermodel.Customer, error)
	FindByID(context.Context, string) (*userdto.CustomerResponse, error)
}

type refreshTokenRepository interface {
	Create(context.Context, *model.RefreshToken) error
	FindByTokenHash(context.Context, string) (*model.RefreshToken, error)
	Revoke(context.Context, string) error
}
