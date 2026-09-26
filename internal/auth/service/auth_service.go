package service

import (
	"context"
	"net/http"
	"strings"
	"time"

	"Dzaakk/simple-commerce/internal/auth/dto"
	"Dzaakk/simple-commerce/internal/auth/model"
	userdto "Dzaakk/simple-commerce/internal/user/dto"
	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/util"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	customers customerService
	tokens    refreshTokenRepository
}

func NewAuthService(customers customerService, tokens refreshTokenRepository) AuthService {
	return &authService{customers: customers, tokens: tokens}
}

func (s *authService) RegisterCustomer(ctx context.Context, req *dto.RegisterCustomerRequest) error {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	fullName := strings.TrimSpace(req.FullName)
	if email == "" || fullName == "" || len(req.Password) < 8 {
		return response.NewAppError(http.StatusBadRequest, "invalid registration data")
	}
	existing, err := s.customers.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existing != nil {
		return response.ErrEmailAlreadyExist
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		return err
	}
	_, err = s.customers.Create(ctx, &userdto.RegisterCustomerRequest{
		Email:        email,
		PasswordHash: hash,
		FullName:     fullName,
	})
	return err
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	customer, err := s.customers.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		return nil, err
	}
	if customer == nil || customer.Status != string(constant.StatusActive) ||
		bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(req.Password)) != nil {
		return nil, response.ErrInvalidCredentials
	}

	accessToken, err := util.GenerateAccessToken(customer.ID, customer.Email)
	if err != nil {
		return nil, response.WrapAppError(http.StatusInternalServerError, "failed to generate access token", err)
	}
	rawRefresh, refreshHash, err := generateRefreshToken()
	if err != nil {
		return nil, response.WrapAppError(http.StatusInternalServerError, "failed to generate refresh token", err)
	}
	if err := s.tokens.Create(ctx, &model.RefreshToken{
		UserID: customer.ID, TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(refreshTokenDuration), CreatedAt: time.Now(),
	}); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: accessToken, RefreshToken: rawRefresh,
		ExpiresIn: int(util.AccessTokenDuration.Seconds()),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, raw string) (*dto.RefreshTokenResponse, error) {
	stored, err := s.tokens.FindByTokenHash(ctx, hashRefreshToken(raw))
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, response.ErrInvalidRefreshToken
	}
	customer, err := s.customers.FindByID(ctx, stored.UserID)
	if err != nil {
		return nil, err
	}
	if customer == nil || customer.Status != string(constant.StatusActive) {
		return nil, response.ErrInvalidRefreshToken
	}
	accessToken, err := util.GenerateAccessToken(customer.ID, customer.Email)
	if err != nil {
		return nil, response.WrapAppError(http.StatusInternalServerError, "failed to generate access token", err)
	}
	return &dto.RefreshTokenResponse{AccessToken: accessToken, ExpiresIn: int(util.AccessTokenDuration.Seconds())}, nil
}

func (s *authService) Logout(ctx context.Context, raw string) error {
	return s.tokens.Revoke(ctx, hashRefreshToken(raw))
}
