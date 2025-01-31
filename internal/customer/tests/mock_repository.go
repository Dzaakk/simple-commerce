package tests

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) Create(ctx context.Context, customer *model.TCustomers) (*int, error) {
	args := m.Called(ctx, customer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockCustomerRepository) FindById(ctx context.Context, id int) (*model.TCustomers, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TCustomers), args.Error(1)
}

func (m *MockCustomerRepository) UpdateBalance(ctx context.Context, id int, balance float64) (*float64, error) {
	args := m.Called(ctx, id, balance)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}

func (m *MockCustomerRepository) GetBalance(ctx context.Context, id int) (*model.CustomerBalance, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CustomerBalance), args.Error(1)
}

func (m *MockCustomerRepository) FindByEmail(ctx context.Context, email string) (*model.TCustomers, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.TCustomers), args.Error(1)
}
