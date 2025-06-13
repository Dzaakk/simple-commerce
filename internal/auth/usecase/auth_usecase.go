package usecase

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"context"
)

type AuthUseCase interface {
	RegistrationCustomer(ctx context.Context, data model.CustomerRegistrationReq) error
	ActivationCustomer(ctx context.Context, data model.CustomerActivationReq) error
	LoginCustomer(ctx context.Context, data model.LoginReq) error
	SellerRegistration(ctx context.Context, data model.SellerRegistrationReq) (*int64, error)
}
