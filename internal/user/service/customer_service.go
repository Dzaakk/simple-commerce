package service

import (
	"Dzaakk/simple-commerce/internal/user/dto"
	"Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"net/http"
	"strconv"
)

type CustomerServiceImpl struct {
	Repo CustomerRepository
}

func NewCustomerService(repo CustomerRepository) CustomerService {
	return &CustomerServiceImpl{Repo: repo}
}

func (c *CustomerServiceImpl) Create(ctx context.Context, req *dto.RegisterCustomerRequest) (string, error) {

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
		return response.NewAppError(http.StatusBadRequest, "invalid parameter customer id")
	}

	if customerID <= 0 {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter customer id")
	}

	data := req.ToUpdateData(customerID)

	rowsAffected, err := c.Repo.Update(ctx, data)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return response.NewAppError(http.StatusNotFound, "customer not found")
	}

	return nil
}

func (c *CustomerServiceImpl) FindByEmail(ctx context.Context, email string) (*model.Customer, error) {

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

func (c *CustomerServiceImpl) UpdateStatus(ctx context.Context, customerID string, status constant.UserStatus) error {
	return c.Repo.UpdateStatus(ctx, customerID, status)
}

func (c *CustomerServiceImpl) UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, customerID string, status constant.UserStatus) error {
	return c.Repo.UpdateStatusWithTx(ctx, tx, customerID, status)
}
