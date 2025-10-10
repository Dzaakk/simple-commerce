package repository

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	"context"
	"database/sql"
)

type CustomerRepository interface {
	Create(ctx context.Context, data model.TCustomers) (int64, error)
	Update(ctx context.Context, data model.TCustomers) (int64, error)
	FindByID(ctx context.Context, customerID int64) (*model.TCustomers, error)
	FindByEmail(ctx context.Context, email string) (*model.TCustomers, error)
	GetBalanceWithTx(ctx context.Context, tx *sql.Tx, customerID int64) (*model.CustomerBalance, error)
	UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, customerID int64, newBalance float64) error
}
