package usecase

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	"Dzaakk/simple-commerce/package/constant"
	"context"
	"errors"
	"strconv"
	"time"
)

type CustomerUsecaseImpl struct {
	Repo CustomerRepository
}

func NewCustomerUseCase(repo CustomerRepository) CustomerUsecase {
	return &CustomerUsecaseImpl{Repo: repo}
}

func (c *CustomerUsecaseImpl) Create(ctx context.Context, req *model.CreateReq) (int64, error) {

	dateOfBirth, err := time.Parse(constant.DateLayout, req.DateOfBirth)
	if err != nil {
		return 0, err
	}

	data := req.ToCreateData(dateOfBirth)

	id, err := c.Repo.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *CustomerUsecaseImpl) Update(ctx context.Context, req *model.UpdateReq) error {

	dateOfBirth, err := time.Parse(constant.DateLayout, req.DateOfBirth)
	if err != nil {
		return err
	}
	customerID, err := strconv.ParseInt(req.CustomerID, 0, 64)
	if err != nil {
		return err
	}

	data := req.ToUpdateData(dateOfBirth, customerID)

	rowsAffected, err := c.Repo.Update(ctx, data)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}

	return nil
}

func (c *CustomerUsecaseImpl) FindByEmail(ctx context.Context, email string) (*model.TCustomers, error) {
	data, err := c.Repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *CustomerUsecaseImpl) FindByID(ctx context.Context, customerID int64) (*model.CustomerRes, error) {
	// validate id
	data, err := c.Repo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	customer := data.ToResponse()

	return &customer, nil
}
