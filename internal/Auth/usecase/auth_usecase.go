package usecase

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"context"
)

type AuthUseCase interface {
	CustomerRegistration(ctx context.Context, data model.CustomerRegistration) (*int64, error)
	SellerRegistration(ctx context.Context, data model.SellerRegistration) (*int64, error)
}
