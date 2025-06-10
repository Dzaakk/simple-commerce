package repository

import (
	"context"
)

type AuthCacheRepository interface {
	SetActivationCustomer(c context.Context, email string, activationCode string) error
	GetActivationCustomer(c context.Context, email string) (string, error)
}
