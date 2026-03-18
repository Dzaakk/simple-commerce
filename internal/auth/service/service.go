package service

import (
	"Dzaakk/simple-commerce/internal/auth/dto"
	"Dzaakk/simple-commerce/internal/auth/model"
	emailmodel "Dzaakk/simple-commerce/internal/email/model"
	userdto "Dzaakk/simple-commerce/internal/user/dto"
	usermodel "Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
	"context"
	"database/sql"
)

type AuthService interface {
	RegisterCustomer(ctx context.Context, req *dto.RegisterCustomerRequest) error
	RegisterSeller(ctx context.Context, req *dto.RegisterSellerRequest) error
	VerifyEmail(ctx context.Context, activationCode string) error
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type customerService interface {
	Create(ctx context.Context, req *userdto.RegisterCustomerRequest) (string, error)
	FindByEmail(ctx context.Context, email string) (*usermodel.Customer, error)
	FindByID(ctx context.Context, customerID string) (*userdto.CustomerRes, error)
	UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, customerID string, status constant.UserStatus) error
}

type sellerService interface {
	Create(ctx context.Context, req *userdto.RegisterSellerRequest) (string, error)
	FindByEmail(ctx context.Context, email string) (*usermodel.Seller, error)
	FindByID(ctx context.Context, sellerID string) (*userdto.SellerRes, error)
	UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, sellerID string, status constant.UserStatus) error
}

type emailService interface {
	SendEmailVerification(ctx context.Context, req emailmodel.VerificationEmailReq) error
}

type activationCodeRepository interface {
	Create(ctx context.Context, data *model.ActivationCode) (int64, error)
	FindByCode(ctx context.Context, code string) (*model.ActivationCode, error)
	MarkAsUsedWithTx(ctx context.Context, tx *sql.Tx, id int64) error
}

type refreshTokenRepository interface {
	Create(ctx context.Context, data *model.RefreshToken) (int64, error)
	FindByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string) error
	RevokeAllByUser(ctx context.Context, userID string, userType constant.UserType) error
}
