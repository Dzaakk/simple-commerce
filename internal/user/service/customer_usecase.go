package usecase

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	"context"
)

type CustomerUsecase interface {
	Create(ctx context.Context, req *model.CreateReq) (int64, error)
	Update(ctx context.Context, req *model.UpdateReq) error
	FindByEmail(ctx context.Context, email string) (*model.TCustomers, error)
	FindByID(ctx context.Context, customerID int64) (*model.CustomerRes, error)
}

type CustomerRepository interface {
	Create(ctx context.Context, data *model.TCustomers) (int64, error)
	Update(ctx context.Context, data *model.TCustomers) (int64, error)
	FindByID(ctx context.Context, customerID int64) (*model.TCustomers, error)
	FindByEmail(ctx context.Context, email string) (*model.TCustomers, error)
}
