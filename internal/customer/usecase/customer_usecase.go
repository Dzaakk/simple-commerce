package usecase

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	"context"
)

type CustomerUseCase interface {
	Update(ctx context.Context, data model.UpdateReq) error
	FindByEmail(ctx context.Context, email string) (*model.TCustomers, error)
	FindByID(ctx context.Context, customerID int64) (*model.CustomerRes, error)
}
