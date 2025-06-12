package repository

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"context"
)

type AuthCacheRepository interface {
	SetActivationCustomer(c context.Context, email string, activationCode string) error
	GetActivationCustomer(c context.Context, email string) (string, error)
	SetTokenCustomer(c context.Context, email, token string) error
	GetTokenCustomer(c context.Context, email string) (*string, error)
	SetCustomerRegistration(c context.Context, data model.CustomerRegistrationReq) error
	GetCustomerRegistration(c context.Context, email string) (*model.CustomerRegistrationReq, error)
}
