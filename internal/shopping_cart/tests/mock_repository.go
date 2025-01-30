package tests

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"

	"github.com/stretchr/testify/mock"
)

type MockShoppingCartRepository struct {
	mock.Mock
}

func (m *MockShoppingCartRepository) FindByID(id int) (*model.TShoppingCart, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.TShoppingCart), args.Error(1)
}
