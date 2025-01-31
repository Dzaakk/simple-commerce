package tests

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockShoppingCartRepository struct {
	mock.Mock
}

func (m *MockShoppingCartRepository) Create(ctx context.Context, shoppingCart *model.TShoppingCart) (*model.TShoppingCart, error) {
	args := m.Called(ctx, shoppingCart)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TShoppingCart), args.Error(1)
}

func (m *MockShoppingCartRepository) FindByID(ctx context.Context, id int) (*model.TShoppingCart, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.TShoppingCart), args.Error(1)
}

func (m *MockShoppingCartRepository) FindByStatusAndCustomerID(ctx context.Context, status string, id int) (*model.TShoppingCart, error) {
	args := m.Called(ctx, status, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.TShoppingCart), args.Error(1)
}
