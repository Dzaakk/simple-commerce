package service

import (
	"context"

	"Dzaakk/simple-commerce/internal/user/dto"
	"Dzaakk/simple-commerce/internal/user/model"
)

type CustomerService interface {
	Create(context.Context, *dto.RegisterCustomerRequest) (string, error)
	FindByEmail(context.Context, string) (*model.Customer, error)
	FindByID(context.Context, string) (*dto.CustomerResponse, error)
}

type CustomerRepository interface {
	Create(context.Context, *model.Customer) (string, error)
	FindByEmail(context.Context, string) (*model.Customer, error)
	FindByID(context.Context, string) (*model.Customer, error)
}
