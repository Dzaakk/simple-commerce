package usecase

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
)

type CustomerUseCase interface {
	Create(data model.CustomerReq) (*int, error)
	FindById(id int) (*model.CustomerRes, error)
	UpdateBalance(id int, balance float64, actionType string) (*float64, error)
	GetBalance(id int) (*model.CustomerBalance, error)
	FindByEmail(email string) (*model.TCustomers, error)
	Update(data model.TCustomers) (*model.CustomerRes, error)
}
