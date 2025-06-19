package usecase

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	eModel "Dzaakk/simple-commerce/internal/email/model"
	"context"
)

type AuthUseCase interface {
	RegistrationCustomer(ctx context.Context, data model.CustomerRegistrationReq) (*eModel.ActivationEmailReq, error)
	ActivationCustomer(ctx context.Context, data model.ActivationReq) error
	LoginCustomer(ctx context.Context, data model.LoginReq) error

	RegistrationSeller(ctx context.Context, data model.SellerRegistrationReq) (*eModel.ActivationEmailReq, error)
	ActivationSeller(ctx context.Context, data model.ActivationReq) error
	LoginSeller(ctx context.Context, data model.LoginReq) error

	Logout(ctx context.Context, email, role string) error
}
