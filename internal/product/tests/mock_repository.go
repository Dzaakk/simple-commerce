package tests

import (
	model "Dzaakk/simple-commerce/internal/product/models"

	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product model.TProduct) (*model.TProduct, error) {
	args := m.Called(product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.TProduct), args.Error(1)
}
