package usecase

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	"context"
)

type CustomerUseCase interface {
	Update(ctx context.Context, data model.UpdateReq) error
	FindByEmail(ctx context.Context, email string) (*model.TCustomers, error)
	FindByID(ctx context.Context, customerID int64) (*model.CustomerRes, error)

	// FindByUsername(ctx context.Context, username string) (*model.DataRes, error)
	// UpdateBalance(ctx context.Context, customerID int64, balance float64, actionType string) (int64, error)
	// GetBalance(ctx context.Context, customerID int64) (*model.CustomerBalanceRes, error)
	// IncreaseBalance(ctx context.Context, customerID int64, amount float64) (*model.CustomerBalanceRes, error)
	// DecreaseBalance(ctx context.Context, customerID int64, amount float64) (*model.CustomerBalanceRes, error)
	// UpdatePassword(ctx context.Context, customerID int64, newPassword string) (int64, error)
	// Deactivate(ctx context.Context, customerID int64) (int64, error)
	//inquiry balance
}
