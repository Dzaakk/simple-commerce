package test

import (
	models "Dzaakk/simple-commerce/internal/customer/models"
	"errors"
)

type MockRepository struct {
	data map[int]models.TCustomers
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		data: make(map[int]models.TCustomers),
	}
}

func (m *MockRepository) Create(data models.TCustomers) (*int, error) {
	id := len(m.data) + 1
	data.Id = id
	m.data[id] = data
	return &id, nil
}

func (m *MockRepository) FindById(id int) (*models.TCustomers, error) {
	if customer, exists := m.data[id]; exists {
		return &customer, nil
	}

	return nil, errors.New("customer not found")
}

func (m *MockRepository) UpdateBalance(id int, balance float32) (*float32, error) {
	if customer, exists := m.data[id]; exists {
		customer.Balance = balance
		m.data[id] = customer
		return &balance, nil
	}

	return nil, errors.New("customer not found")
}

func (m *MockRepository) GetBalance(id int) (*models.CustomerBalance, error) {
	if customer, exists := m.data[id]; exists {
		balance := models.CustomerBalance{
			Id:      id,
			Balance: customer.Balance,
		}
		return &balance, nil
	}

	return nil, errors.New("customer not found")
}

func (m *MockRepository) FindByEmail(email string) (*models.TCustomers, error) {
	for _, customer := range m.data {
		if customer.Email == email {
			return &customer, nil

		}
	}

	return nil, errors.New("customer not found")
}
