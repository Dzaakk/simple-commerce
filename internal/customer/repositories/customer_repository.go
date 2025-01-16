package repositories

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	"context"
	"database/sql"
)

type CustomerRepository interface {
	Create(ctx context.Context, data model.TCustomers) (int64, error)
	FindById(ctx context.Context, id int64) (*model.TCustomers, error)
	UpdateBalance(ctx context.Context, id int64, balance float64) (int64, error)
	GetBalance(ctx context.Context, id int64) (*model.CustomerBalance, error)
	InquiryBalance(ctx context.Context, id int64) (float64, error)
	FindByEmail(ctx context.Context, email string) (*model.TCustomers, error)
	GetBalanceWithTx(ctx context.Context, tx *sql.Tx, id int64) (*model.CustomerBalance, error)
	UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, id int64, newBalance float64) error
	UpdatePassword(ctx context.Context, id int64, newPassword string) (int64, error)
	Deactive(ctx context.Context, id int64) (int64, error)
}
