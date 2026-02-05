package usecase

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	eModel "Dzaakk/simple-commerce/internal/email/model"
	"context"
)

type AuthUsecase interface {
	RegistrationCustomer(ctx context.Context, data model.CustomerRegistrationReq) (*eModel.ActivationEmailReq, error)
	ActivationCustomer(ctx context.Context, data model.ActivationReq) error
	LoginCustomer(ctx context.Context, data model.LoginReq) error

	RegistrationSeller(ctx context.Context, data model.SellerRegistrationReq) (*eModel.ActivationEmailReq, error)
	ActivationSeller(ctx context.Context, data model.ActivationReq) error
	LoginSeller(ctx context.Context, data model.LoginReq) error

	Logout(ctx context.Context, email, role string) error
}

type AuthCacheCustomer interface {
	SetActivation(ctx context.Context, email, activationCode string) error
	GetActivation(ctx context.Context, email string) (string, error)
	SetToken(ctx context.Context, email, token string) error
	GetToken(ctx context.Context, email string) (*string, error)
	SetRegistration(ctx context.Context, data model.CustomerRegistrationReq) error
	GetRegistration(ctx context.Context, email string) (*model.CustomerRegistrationReq, error)
	DeleteToken(ctx context.Context, email string) error
}

type AuthCacheSeller interface {
	SetActivation(ctx context.Context, email string, activationCode string) error
	GetActivation(ctx context.Context, email string) (string, error)
	SetToken(ctx context.Context, email, token string) error
	GetToken(ctx context.Context, email string) (*string, error)
	SetRegistration(ctx context.Context, data model.SellerRegistrationReq) error
	GetRegistration(ctx context.Context, email string) (*model.SellerRegistrationReq, error)
	DeleteToken(ctx context.Context, email string) error
}

type CustomerTokenUsecase interface {
	GetToken(ctx context.Context, email string) (*string, error)
}

type SellerTokenUsecase interface {
	GetToken(ctx context.Context, email string) (*string, error)
}
