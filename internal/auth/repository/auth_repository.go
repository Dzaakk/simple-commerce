package repository

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"context"
)

type AuthCacheCustomer interface {
	SetActivationCustomer(ctx context.Context, email string, activationCode string) error
	GetActivationCustomer(ctx context.Context, email string) (string, error)
	SetTokenCustomer(ctx context.Context, email, token string) error
	GetTokenCustomer(ctx context.Context, email string) (*string, error)
	SetCustomerRegistration(ctx context.Context, data model.CustomerRegistrationReq) error
	GetCustomerRegistration(ctx context.Context, email string) (*model.CustomerRegistrationReq, error)
}

type AuthCacheSeller interface {
	SetActivationSeller(ctx context.Context, email string, activationCode string) error
	GetActivationSeller(ctx context.Context, email string) (string, error)
	SetTokenSeller(ctx context.Context, email, token string) error
	GetTokenSeller(ctx context.Context, email string) (*string, error)
	SetSellerRegistration(ctx context.Context, data model.SellerRegistrationReq) error
	GetSellerRegistration(ctx context.Context, email string) (*model.SellerRegistrationReq, error)
}
