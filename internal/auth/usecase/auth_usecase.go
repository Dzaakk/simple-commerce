package usecase

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"context"
)

type AuthUseCase interface {
	RegistrationCustomer(ctx context.Context, data model.CustomerRegistrationReq) error
	ActivationCustomer(ctx context.Context, data model.ActivationReq) error
	LoginCustomer(ctx context.Context, data model.LoginReq) error

	RegistrationSeller(ctx context.Context, data model.SellerRegistrationReq) error
	ActivationSeller(ctx context.Context, data model.ActivationReq) error
	LoginSeller(ctx context.Context, data model.LoginReq) error
}
