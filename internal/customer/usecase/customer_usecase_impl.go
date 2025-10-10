package usecase

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	repo "Dzaakk/simple-commerce/internal/customer/repository"
	"context"
	"strconv"
	"time"
)

type CustomerUseCaseImpl struct {
	Repo repo.CustomerRepository
}

func NewCustomerUseCase(repo repo.CustomerRepository) CustomerUseCase {
	return &CustomerUseCaseImpl{Repo: repo}
}

func (c *CustomerUseCaseImpl) Update(ctx context.Context, req model.UpdateReq) error {

	dateOfBirth, err := time.Parse("02-01-2006", req.DateOfBirth)
	if err != nil {
		return err
	}
	customerID, err := strconv.ParseInt(req.CustomerID, 0, 64)
	if err != nil {
		return err
	}

	data := req.ToCustomerModel(dateOfBirth, customerID)

	rowsAffected, err := c.Repo.Update(ctx, data)
	if err != nil || rowsAffected == 0 {
		return err
	}

	return nil
}

func (c *CustomerUseCaseImpl) FindByEmail(ctx context.Context, email string) (*model.TCustomers, error) {
	data, err := c.Repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *CustomerUseCaseImpl) FindByID(ctx context.Context, customerID int64) (*model.CustomerRes, error) {
	data, err := c.Repo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	customer := data.ToResponse()

	return &customer, nil
}
