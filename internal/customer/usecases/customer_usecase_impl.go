package usecases

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	repo "Dzaakk/simple-commerce/internal/customer/repositories"
	template "Dzaakk/simple-commerce/package/templates"
	"fmt"
	"time"
)

type CustomerUseCaseImpl struct {
	repo repo.CustomerRepository
}

func NewCustomerUseCase(repo repo.CustomerRepository) CustomerUseCase {
	return &CustomerUseCaseImpl{repo}
}

func (c *CustomerUseCaseImpl) Create(data model.CreateReq) (int64, error) {
	hashedPassword, err := template.HashPassword(data.Password)
	if err != nil {
		return 0, err
	}

	customer := model.TCustomers{
		Username:    data.Username,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    string(hashedPassword),
		Balance:     float64(10000000),
		Status:      "A",
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: "system",
		},
	}

	customerId, err := c.repo.Create(customer)
	if err != nil {
		return 0, err
	}
	return customerId, nil
}

func (c *CustomerUseCaseImpl) FindByEmail(email string) (*model.TCustomers, error) {
	data, err := c.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *CustomerUseCaseImpl) FindById(id int64) (*model.DataRes, error) {
	data, err := c.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	customer := &model.DataRes{
		Id:          fmt.Sprintf("%v", data.Id),
		Username:    data.Username,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Balance:     fmt.Sprintf("%0.f", data.Balance),
	}

	return customer, nil

}

func (c *CustomerUseCaseImpl) GetBalance(id int64) (*model.CustomerBalanceRes, error) {
	data, err := c.repo.GetBalance(id)
	if err != nil {
		return nil, err
	}

	customer := &model.CustomerBalanceRes{
		Id:      fmt.Sprintf("%d", data.Id),
		Balance: fmt.Sprintf("%.2f", data.Balance),
	}

	return customer, nil
}

func (c *CustomerUseCaseImpl) UpdateBalance(id int64, balance float64, actionType string) (float64, error) {
	data, err := c.repo.GetBalance(id)
	if err != nil {
		return 0, err
	}

	//Add Balance
	if actionType == "A" {
		balance += data.Balance
		balance, err := c.repo.UpdateBalance(id, balance)
		if err != nil {
			return 0, err
		}

		return balance, nil
	}

	//payment
	if actionType == "P" {
		balance = data.Balance - balance
		balance, err := c.repo.UpdateBalance(id, balance)
		if err != nil {
			return 0, err
		}

		return balance, nil
	}

	return 0, nil
}

func (c *CustomerUseCaseImpl) DecreaseBalance(id int64, amount float64) (*model.CustomerBalanceRes, error) {
	if amount < 0 {
		return nil, fmt.Errorf("invalid amount")
	}

	data, err := c.repo.GetBalance(id)
	if err != nil {
		return nil, err
	}

	amount += data.Balance
	balance, err := c.repo.UpdateBalance(id, amount)
	if err != nil {
		return nil, err
	}
	res := &model.CustomerBalanceRes{
		Id:      fmt.Sprintf("%d", id),
		Balance: fmt.Sprintf("%.2f", balance),
	}
	return res, nil
}

func (c *CustomerUseCaseImpl) IncreaseBalance(id int64, amount float64) (*model.CustomerBalanceRes, error) {
	if amount < 0 {
		return nil, fmt.Errorf("invalid amount")
	}

	data, err := c.repo.GetBalance(id)
	if err != nil {
		return nil, err
	}

	data.Balance -= amount
	balance, err := c.repo.UpdateBalance(id, amount)
	if err != nil {
		return nil, err
	}
	res := &model.CustomerBalanceRes{
		Id:      fmt.Sprintf("%d", id),
		Balance: fmt.Sprintf("%.2f", balance),
	}
	return res, nil
}

func (c *CustomerUseCaseImpl) Deactivate(id int64) (int64, error) {

	rowsAffected, err := c.repo.Deactive(id)
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (c *CustomerUseCaseImpl) UpdatePassword(id int64, newPassword string) (int64, error) {
	hashedPassword, err := template.HashPassword(newPassword)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := c.repo.UpdatePassword(id, string(hashedPassword))
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (c *CustomerUseCaseImpl) Update(dataReq model.UpdateReq) (int64, error) {
	panic("unimplemented")
}
