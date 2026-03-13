package service

import (
	"Dzaakk/simple-commerce/internal/auth/dto"
	"Dzaakk/simple-commerce/internal/auth/model"
	usermodel "Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
	"context"
)

type AuthService interface {
	RegisterCustomer(ctx context.Context, req *dto.RegisterCustomerRequest) error
	RegisterSeller(ctx context.Context, req *dto.RegisterSellerRequest) error
	VerifyEmail(ctx context.Context, email, code, userType constant.UserType) error
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type customerService interface {
	Create(ctx context.Context, req *dto.RegisterCustomerRequest) (string, error)
	FindByEmail(ctx context.Context, email string) (*usermodel.Customer, error)
	FindByID(ctx context.Context, customerID string) (*usermodel.Customer, error)
	UpdateStatus(ctx context.Context, customerID string, status constant.UserStatus) error
}

type sellerService interface {
	Create(ctx context.Context, req *dto.RegisterSellerRequest) (string, error)
	FindByEmail(ctx context.Context, email string) (*usermodel.Seller, error)
	FindByID(ctx context.Context, sellerID string) (*usermodel.Seller, error)
	UpdateStatus(ctx context.Context, sellerID string, status constant.UserStatus) error
}

type activationCodeRepository interface {
	Create(ctx context.Context, data *model.ActivationCode) error
	FindByEmailAndCode(ctx context.Context, email, code string) (*model.ActivationCode, error)
	MarkAsUsed(ctx context.Context, id int64) error
}

type refreshTokenRepository interface {
	Create(ctx context.Context, data *model.RefreshToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string) error
	RevokeAllByUser(ctx context.Context, userID string, userType constant.UserType) error
}
