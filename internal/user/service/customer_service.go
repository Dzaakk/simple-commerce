package service

import (
	"Dzaakk/simple-commerce/internal/user/domain"
	"Dzaakk/simple-commerce/internal/user/dto"
	"context"
	"errors"
	"strconv"
)

type CustomerServiceImpl struct {
	Repo CustomerRepository
}

func NewCustomerService(repo CustomerRepository) CustomerService {
	return &CustomerServiceImpl{Repo: repo}
}

func (c *CustomerServiceImpl) Create(ctx context.Context, req *dto.CreateReq) (string, error) {

	data := req.ToCreateData()

	id, err := c.Repo.Create(ctx, data)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (c *CustomerServiceImpl) Update(ctx context.Context, req *dto.UpdateReq) error {

	customerID, err := strconv.ParseInt(req.CustomerID, 0, 64)
	if err != nil {
		return err
	}

	if customerID <= 0 {
		return errors.New("invalid parameter customer id")
	}

	data := req.ToUpdateData(customerID)

	rowsAffected, err := c.Repo.Update(ctx, data)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}

	return nil
}

func (c *CustomerServiceImpl) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {

	data, err := c.Repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *CustomerServiceImpl) FindByID(ctx context.Context, customerID string) (*dto.CustomerRes, error) {

	data, err := c.Repo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	customer := dto.ToCustomerRes(data)

	return &customer, nil
}
