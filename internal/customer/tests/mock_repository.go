package tests

import (
	model "Dzaakk/simple-commerce/internal/customer/models"

	"github.com/stretchr/testify/mock"
)

type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) Create(customer *model.TCustomers) (*int, error) {
	args := m.Called(customer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockCustomerRepository) FindById(id int) (*model.TCustomers, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TCustomers), args.Error(1)
}

func (m *MockCustomerRepository) UpdateBalance(id int, balance float64) (*float64, error) {
	args := m.Called(id, balance)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*float64), args.Error(1)
}
