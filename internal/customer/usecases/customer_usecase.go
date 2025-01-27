package usecases

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	"context"
)

type CustomerUseCase interface {
	Create(ctx context.Context, data model.CreateReq) (int64, error)
	FindById(ctx context.Context, id int64) (*model.DataRes, error)
	UpdateBalance(ctx context.Context, id int64, balance float64, actionType string) (int64, error)
	GetBalance(ctx context.Context, id int64) (*model.CustomerBalanceRes, error)
	IncreaseBalance(ctx context.Context, id int64, amount float64) (*model.CustomerBalanceRes, error)
	DecreaseBalance(ctx context.Context, id int64, amount float64) (*model.CustomerBalanceRes, error)
	FindByEmail(ctx context.Context, email string) (*model.TCustomers, error)
	Update(ctx context.Context, data model.UpdateReq) (int64, error)
	UpdatePassword(ctx context.Context, id int64, newPassword string) (int64, error)
	Deactivate(ctx context.Context, id int64) (int64, error)
	//inquiry balance
}
