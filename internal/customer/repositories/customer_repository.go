package repositories

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	"database/sql"
)

type CustomerRepository interface {
	Create(data model.TCustomers) (int64, error)
	FindById(id int64) (*model.TCustomers, error)
	UpdateBalance(id int64, balance float64) (float64, error)
	GetBalance(id int64) (*model.CustomerBalance, error)
	FindByEmail(email string) (*model.TCustomers, error)
	GetBalanceWithTx(tx *sql.Tx, id int64) (*model.CustomerBalance, error)
	UpdateBalanceWithTx(tx *sql.Tx, id int64, newBalance float64) error
}
