package service

import (
	"Dzaakk/simple-commerce/internal/auth/dto"
	"Dzaakk/simple-commerce/internal/auth/model"
	emailmodel "Dzaakk/simple-commerce/internal/email/model"
	userdto "Dzaakk/simple-commerce/internal/user/dto"
	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type authService struct {
	db             *sql.DB
	customerSvc    customerService
	sellerSvc      sellerService
	emailService   emailService
	activationRepo activationCodeRepository
	refreshRepo    refreshTokenRepository
}

func NewAuthService(
	db *sql.DB,
	customerSvc customerService,
	sellerSvc sellerService,
	emailService emailService,
	activationRepo activationCodeRepository,
	refreshRepo refreshTokenRepository,
) AuthService {
	return &authService{
		db:             db,
		customerSvc:    customerSvc,
		sellerSvc:      sellerSvc,
		emailService:   emailService,
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

	createReq := &userdto.RegisterCustomerRequest{
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		Phone:    req.Phone,
	}

	_, err = s.customerSvc.Create(ctx, createReq)
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

	_, err = s.activationRepo.Create(ctx, activationData)
	if err != nil {
		return err
	}

	// send activation email
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recovered from panic in email goroutine: %v", r)
			}
		}()
		err := s.emailService.SendEmailVerification(context.Background(), emailmodel.VerificationEmailReq{
			Email:          req.Email,
			Username:       req.FullName,
			ActivationLink: fmt.Sprintf("http://localhost:8080/api/v1/auth/verify-email?code=%s", activationCode),
		})
		if err != nil {
			log.Printf("failed to send email to %s: %v", req.Email, err)
		}
	}()

	return nil
}

func (s *authService) RegisterSeller(ctx context.Context, req *dto.RegisterSellerRequest) error {

	// check if email already exist
	existingSeller, err := s.sellerSvc.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existingSeller != nil {
		return response.ErrEmailAlreadyExist
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return err
	}

	createReq := &userdto.RegisterSellerRequest{
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		Phone:    req.Phone,
		ShopName: req.ShopName,
	}

	_, err = s.sellerSvc.Create(ctx, createReq)
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
		UserType:  string(constant.Seller),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	_, err = s.activationRepo.Create(ctx, activationData)
	if err != nil {
		return err
	}

	// send activation email

	return nil
}

func (s *authService) VerifyEmail(ctx context.Context, activationCode string) error {
	activationData, err := s.activationRepo.FindByCode(ctx, activationCode)
	if err != nil {
		return err
	}
	if activationData == nil {
		return response.ErrInvalidActivationCode
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return response.Error("failed to begin transaction", err)
	}
	defer tx.Rollback()

	switch activationData.UserType {
	case string(constant.Customer):
		customer, err := s.customerSvc.FindByEmail(ctx, activationData.Email)
		if err != nil {
			return err
		}
		if customer == nil {
			return response.ErrUserNotFound
		}

		err = s.customerSvc.UpdateStatusWithTx(ctx, tx, customer.ID, constant.StatusActive)
		if err != nil {
			return err
		}
	case string(constant.Seller):
		seller, err := s.sellerSvc.FindByEmail(ctx, activationData.Email)
		if err != nil {
			return err
		}
		if seller == nil {
			return response.ErrUserNotFound
		}

		err = s.sellerSvc.UpdateStatusWithTx(ctx, tx, seller.ID, constant.StatusActive)
		if err != nil {
			return err
		}
	default:
		return response.ErrInvalidActivationCode
	}

	err = s.activationRepo.MarkAsUsedWithTx(ctx, tx, activationData.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
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
