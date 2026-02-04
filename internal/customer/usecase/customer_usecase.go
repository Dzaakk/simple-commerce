package usecase

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	"context"
)

type CustomerRepository interface {
	Create(ctx context.Context, data model.TCustomers) (int64, error)
	Update(ctx context.Context, data model.TCustomers) (int64, error)
	FindByID(ctx context.Context, customerID int64) (*model.TCustomers, error)
	FindByEmail(ctx context.Context, email string) (*model.TCustomers, error)
}

type CustomerUseCase interface {
	Update(ctx context.Context, data model.UpdateReq) error
	FindByEmail(ctx context.Context, email string) (*model.TCustomers, error)
	FindByID(ctx context.Context, customerID int64) (*model.CustomerRes, error)
}
