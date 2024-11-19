package repositories

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	"database/sql"
)

type CustomerRepository interface {
	Create(data model.TCustomers) (*int, error)
	FindById(id int) (*model.TCustomers, error)
	UpdateBalance(id int, balance float64) (*float64, error)
	GetBalance(id int) (*model.CustomerBalance, error)
	FindByEmail(email string) (*model.TCustomers, error)
	GetBalanceWithTx(tx *sql.Tx, id int) (*model.CustomerBalance, error)
	UpdateBalanceWithTx(tx *sql.Tx, id int, newBalance float64) error
}
