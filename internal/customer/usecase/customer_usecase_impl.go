package usecase

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	repo "Dzaakk/simple-commerce/internal/customer/repository"
	"Dzaakk/simple-commerce/package/util"
	"context"
	"fmt"
)

type CustomerUseCaseImpl struct {
	repo repo.CustomerRepository
}

func NewCustomerUseCase(repo repo.CustomerRepository) CustomerUseCase {
	return &CustomerUseCaseImpl{repo}
}

func (c *CustomerUseCaseImpl) FindByEmail(ctx context.Context, email string) (*model.TCustomers, error) {
	data, err := c.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *CustomerUseCaseImpl) FindByID(ctx context.Context, customerID int64) (*model.DataRes, error) {
	data, err := c.repo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	customer := &model.DataRes{
		CustomerID:  fmt.Sprintf("%v", data.ID),
		Username:    data.Username,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Balance:     fmt.Sprintf("%0.f", data.Balance),
	}

	return customer, nil

}

func (c *CustomerUseCaseImpl) GetBalance(ctx context.Context, customerID int64) (*model.CustomerBalanceRes, error) {
	data, err := c.repo.GetBalance(ctx, customerID)
	if err != nil {
		return nil, err
	}

	customer := &model.CustomerBalanceRes{
		CustomerID: fmt.Sprintf("%d", data.CustomerID),
		Balance:    fmt.Sprintf("%.2f", data.Balance),
	}

	return customer, nil
}

func (c *CustomerUseCaseImpl) UpdateBalance(ctx context.Context, customerID int64, balance float64, actionType string) (int64, error) {
	data, err := c.repo.GetBalance(ctx, customerID)
	if err != nil {
		return 0, err
	}

	//Add Balance
	if actionType == "A" {
		balance += data.Balance
		balance, err := c.repo.UpdateBalance(ctx, customerID, balance)
		if err != nil {
			return 0, err
		}

		return balance, nil
	}

	//payment
	if actionType == "P" {
		balance = data.Balance - balance
		balance, err := c.repo.UpdateBalance(ctx, customerID, balance)
		if err != nil {
			return 0, err
		}

		return balance, nil
	}

	return 0, nil
}

func (c *CustomerUseCaseImpl) DecreaseBalance(ctx context.Context, customerID int64, amount float64) (*model.CustomerBalanceRes, error) {
	if amount < 0 {
		return nil, fmt.Errorf("invalid amount")
	}

	data, err := c.repo.GetBalance(ctx, customerID)
	if err != nil {
		return nil, err
	}

	amount += data.Balance
	balance, err := c.repo.UpdateBalance(ctx, customerID, amount)
	if err != nil {
		return nil, err
	}
	res := &model.CustomerBalanceRes{
		CustomerID: fmt.Sprintf("%d", customerID),
		Balance:    fmt.Sprintf("%d", balance),
	}
	return res, nil
}

func (c *CustomerUseCaseImpl) IncreaseBalance(ctx context.Context, customerID int64, amount float64) (*model.CustomerBalanceRes, error) {
	if amount < 0 {
		return nil, fmt.Errorf("invalid amount")
	}

	data, err := c.repo.GetBalance(ctx, customerID)
	if err != nil {
		return nil, err
	}

	data.Balance -= amount
	balance, err := c.repo.UpdateBalance(ctx, customerID, amount)
	if err != nil {
		return nil, err
	}
	res := &model.CustomerBalanceRes{
		CustomerID: fmt.Sprintf("%d", customerID),
		Balance:    fmt.Sprintf("%d", balance),
	}
	return res, nil
}

func (c *CustomerUseCaseImpl) Deactivate(ctx context.Context, customerID int64) (int64, error) {

	rowsAffected, err := c.repo.Deactive(ctx, customerID)
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (c *CustomerUseCaseImpl) UpdatePassword(ctx context.Context, customerID int64, newPassword string) (int64, error) {
	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := c.repo.UpdatePassword(ctx, customerID, string(hashedPassword))
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (c *CustomerUseCaseImpl) Update(ctx context.Context, dataReq model.UpdateReq) (int64, error) {
	panic("unimplemented")
}
