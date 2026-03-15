package service

import (
	"Dzaakk/simple-commerce/internal/auth/dto"
	"Dzaakk/simple-commerce/internal/auth/model"
	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"time"
)

type authService struct {
	customerSvc    customerService
	sellerSvc      sellerService
	activationRepo activationCodeRepository
	refreshRepo    refreshTokenRepository
}

func NewAuthService(
	customerSvc customerService,
	sellerSvc sellerService,
	activationRepo activationCodeRepository,
	refreshRepo refreshTokenRepository,
) AuthService {
	return &authService{
		customerSvc:    customerSvc,
		sellerSvc:      sellerSvc,
		activationRepo: activationRepo,
		refreshRepo:    refreshRepo,
	}
}

func (s *authService) RegisterCustomer(ctx context.Context, req *dto.RegisterCustomerRequest) error {

	// check if email already exist
	existingCustomer, err := s.customerSvc.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existingCustomer != nil {
		return response.ErrEmailAlreadyExist
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return err
	}

	req.Password = hashedPassword

	_, err = s.customerSvc.Create(ctx, req)
	if err != nil {
		return err
	}

	activationCode, err := generateActivationCode()
	if err != nil {
		return err
	}

	activationData := &model.ActivationCode{
		Email:     req.Email,
		Code:      activationCode,
		UserType:  string(constant.Customer),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	err = s.activationRepo.Create(ctx, activationData)
	if err != nil {
		return err
	}

	// send activation email

	return nil
}

func (s *authService) RegisterSeller(ctx context.Context, req *dto.RegisterSellerRequest) error {
	panic("unimplemented")
}

func (s *authService) VerifyEmail(ctx context.Context, email constant.UserType, code constant.UserType, userType constant.UserType) error {
	panic("unimplemented")
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	panic("unimplemented")
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error) {
	panic("unimplemented")
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	panic("unimplemented")
}
