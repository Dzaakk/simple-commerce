package usecase

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
)

type CustomerUseCase interface {
	Create(data model.CustomerReq) (*int, error)
	FindById(id int) (*model.CustomerRes, error)
	UpdateBalance(id int, balance float32, actionType string) (*float32, error)
	GetBalance(id int) (*model.CustomerBalance, error)
	FindByEmail(email string) (*model.TCustomers, error)
	Update(data model.TCustomers) (*model.CustomerRes, error)
}
