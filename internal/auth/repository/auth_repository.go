package repository

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"context"
)

type AuthCacheRepository interface {
	SetActivationCustomer(c context.Context, email string, activationCode string) error
	GetActivationCustomer(c context.Context, email string) (string, error)
	SetTokenCustomer(c context.Context, data model.CustomerToken) error
	GetTokenCustomer(c context.Context, email string) (*model.CustomerToken, error)
}
