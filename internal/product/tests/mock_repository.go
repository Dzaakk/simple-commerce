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

func (m *MockProductRepository) Update(product model.TProduct) (int64, error) {
	args := m.Called(product)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProductRepository) SetStockByProductID(productID, stock int) (int64, error) {
	args := m.Called(productID, stock)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProductRepository) FindByProductID(productID int) (*model.TProduct, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.TProduct), args.Error(1)
}

func (m *MockProductRepository) FindByCategoryID(categoryID int) ([]*model.TProduct, error) {
	args := m.Called(categoryID)
	return args.Get(0).([]*model.TProduct), args.Error(1)
}

func (m *MockProductRepository) FindByProductName(productName string) (*model.TProduct, error) {
	args := m.Called(productName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.TProduct), args.Error(1)
}

func (m *MockProductRepository) FindBySellerID(sellerID int) ([]*model.TProduct, error) {
	args := m.Called(sellerID)
	return args.Get(0).([]*model.TProduct), args.Error(1)
}

func (m *MockProductRepository) FindBySellerIDAndCategoryID(sellerID, categoryID int) ([]*model.TProduct, error) {
	args := m.Called(sellerID, categoryID)
	return args.Get(0).([]*model.TProduct), args.Error(1)
}
