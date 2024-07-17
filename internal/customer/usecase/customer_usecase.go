package customer

import (
	model "Dzaakk/synapsis/internal/customer/models"
	repo "Dzaakk/synapsis/internal/customer/repository"
	"Dzaakk/synapsis/package/template"
	"fmt"
	"time"
)

type CustomerUseCase interface {
	Create(data model.CustomerReq) (*int, error)
	FindById(id int) (*model.CustomerRes, error)
	UpdateBalance(id int, balance float32, actionType string) (*float32, error)
	GetBalance(id int) (*model.CustomerBalance, error)
}

type CustomerUseCaseImpl struct {
	repo repo.CustomerRepository
}

func NewCustomerUseCase(repo repo.CustomerRepository) CustomerUseCase {
	return &CustomerUseCaseImpl{repo}
}

func (c *CustomerUseCaseImpl) Create(data model.CustomerReq) (*int, error) {
	customer := model.TCustomers{
		Username:    data.Username,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    data.Password,
		Balance:     float32(1000000000),
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: "system",
		},
	}

	customerId, err := c.repo.Create(customer)
	if err != nil {
		return nil, err
	}
	return customerId, nil
}

func (c *CustomerUseCaseImpl) FindById(id int) (*model.CustomerRes, error) {
	data, err := c.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	customer := &model.CustomerRes{
		Id:          fmt.Sprintf("%v", data.Id),
		Username:    data.Username,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Balance:     fmt.Sprintf("%v", data.Balance),
	}

	return customer, nil

}

func (c *CustomerUseCaseImpl) GetBalance(id int) (*model.CustomerBalance, error) {
	data, err := c.repo.GetBalance(id)
	if err != nil {
		return nil, err
	}

	customer := &model.CustomerBalance{
		Id:      data.Id,
		Balance: data.Balance,
	}

	return customer, nil
}

func (c *CustomerUseCaseImpl) UpdateBalance(id int, balance float32, actionType string) (*float32, error) {
	data, err := c.repo.GetBalance(id)
	if err != nil {
		return nil, err
	}

	//Add Balance
	if actionType == "A" {
		balance += data.Balance
		balance, err := c.repo.UpdateBalance(id, balance)
		if err != nil {
			return nil, err
		}

		return balance, nil
	}

	//payment
	if actionType == "P" {
		balance = data.Balance - balance
		balance, err := c.repo.UpdateBalance(id, balance)
		if err != nil {
			return nil, err
		}

		return balance, nil
	}

	return nil, nil
}
