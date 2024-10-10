package repository

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
)

type CustomerRepository interface {
	Create(data model.TCustomers) (*int, error)
	FindById(id int) (*model.TCustomers, error)
	UpdateBalance(id int, balance float32) (*float32, error)
	GetBalance(id int) (*model.CustomerBalance, error)
	FindByEmail(email string) (*model.TCustomers, error)
}
