package tests

import (
	model "Dzaakk/simple-commerce/internal/product/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockRepo    = new(MockProductRepository)
	testProduct = model.TProduct{
		Id:          1,
		ProductName: "Monitor",
		Price:       1700000,
		Stock:       10,
		CategoryId:  3,
		SellerId:    1,
	}
	testProductID   int
	testProductName string
)

func TestCreateProduct(t *testing.T) {
	mockRepo.On("Create", testProduct).Return(&testProduct, nil)

	createdProduct, err := mockRepo.Create(testProduct)

	assert.NoError(t, err)
	assert.NotNil(t, createdProduct)

	assert.Equal(t, testProduct.Id, createdProduct.Id)
	assert.Equal(t, testProduct.ProductName, createdProduct.ProductName)
	assert.Equal(t, testProduct.Price, createdProduct.Price)
	assert.Equal(t, testProduct.Stock, createdProduct.Stock)
	assert.Equal(t, testProduct.CategoryId, createdProduct.CategoryId)
	assert.Equal(t, testProduct.SellerId, createdProduct.SellerId)

	mockRepo.AssertExpectations(t)
}

func TestFindByCategoryId(t *testing.T) {
	testProductID = 1

	mockRepo.On("FindById", testProductID).Return(&testProduct, nil)

	foundProduct, err := mockRepo.FindById(testProductID)

	assert.NoError(t, err)
	assert.NotNil(t, foundProduct)

	assert.Equal(t, testProduct.Id, foundProduct.Id)
	assert.Equal(t, testProduct.ProductName, foundProduct.ProductName)
	assert.Equal(t, testProduct.Price, foundProduct.Price)
	assert.Equal(t, testProduct.Stock, foundProduct.Stock)
	assert.Equal(t, testProduct.CategoryId, foundProduct.CategoryId)
	assert.Equal(t, testProduct.SellerId, foundProduct.SellerId)

	mockRepo.AssertExpectations(t)
}
func TestFindByName(t *testing.T) {
	testProductName = "Monitor"
	mockRepo.On("FindByName", testProductName).Return(&testProduct, nil)

	foundProduct, err := mockRepo.FindByName(testProductName)

	assert.NoError(t, err)
	assert.NotNil(t, foundProduct)

	assert.Equal(t, testProduct.Id, foundProduct.Id)
	assert.Equal(t, testProduct.ProductName, foundProduct.ProductName)
	assert.Equal(t, testProduct.Price, foundProduct.Price)
	assert.Equal(t, testProduct.Stock, foundProduct.Stock)
	assert.Equal(t, testProduct.CategoryId, foundProduct.CategoryId)
	assert.Equal(t, testProduct.SellerId, foundProduct.SellerId)

	mockRepo.AssertExpectations(t)
}
