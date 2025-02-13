package usecases

import (
	model "Dzaakk/simple-commerce/internal/auth/models"
	"context"
)

type AuthUseCase interface {
	CustomerRegistration(ctx context.Context, data model.CustomerRegistration) (*int64, error)
	CustomerLogin(ctx context.Context, data model.LoginReq) error
	SellerRegistration(ctx context.Context, data model.SellerRegistration) (*int64, error)
	SellerLogin(ctx context.Context, data model.LoginReq) error
}
