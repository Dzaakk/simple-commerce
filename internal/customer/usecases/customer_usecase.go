package usecases

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
)

type CustomerUseCase interface {
	Create(data model.CreateReq) (int64, error)
	FindById(id int64) (*model.DataRes, error)
	UpdateBalance(id int64, balance float64, actionType string) (float64, error)
	GetBalance(id int64) (*model.CustomerBalance, error)
	FindByEmail(email string) (*model.TCustomers, error)
	Update(data model.TCustomers) (int64, error)
	UpdatePassword(id int64, newPassword string) (int64, error)
	Deactivate(id int64) (int64, error)
}
