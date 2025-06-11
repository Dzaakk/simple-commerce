package usecase

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"context"
)

type AuthUseCase interface {
	CustomerRegistration(ctx context.Context, data model.CustomerRegistrationReq) error
	SellerRegistration(ctx context.Context, data model.SellerRegistrationReq) (*int64, error)
	CustomerActivation(ctx context.Context, data model.CustomerActivationReq) error
	CustomerLogin(ctx context.Context, data model.LoginReq) error
}
