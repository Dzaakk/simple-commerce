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
	testProductID       int
	testProductName     string
	testProductSellerID int
)

func assertProductEquality(t *testing.T, expected, actual *model.TProduct) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.ProductName, actual.ProductName)
	assert.Equal(t, expected.Price, actual.Price)
	assert.Equal(t, expected.Stock, actual.Stock)
	assert.Equal(t, expected.CategoryId, actual.CategoryId)
	assert.Equal(t, expected.SellerId, actual.SellerId)
}

func TestCreateProduct(t *testing.T) {
	mockRepo.On("Create", testProduct).Return(&testProduct, nil)

	createdProduct, err := mockRepo.Create(testProduct)

	assert.NoError(t, err)
	assert.NotNil(t, createdProduct)
	assertProductEquality(t, &testProduct, createdProduct)
	mockRepo.AssertExpectations(t)
}

func TestFindByCategoryId(t *testing.T) {
	testProductID = 1

	mockRepo.On("FindById", testProductID).Return(&testProduct, nil)

	foundProduct, err := mockRepo.FindById(testProductID)

	assert.NoError(t, err)
	assert.NotNil(t, foundProduct)

	assertProductEquality(t, &testProduct, foundProduct)

	mockRepo.AssertExpectations(t)
}

func TestFindByName(t *testing.T) {
	testProductName = "Monitor"
	mockRepo.On("FindByName", testProductName).Return(&testProduct, nil)

	foundProduct, err := mockRepo.FindByName(testProductName)

	assert.NoError(t, err)
	assert.NotNil(t, foundProduct)

	assertProductEquality(t, &testProduct, foundProduct)

	mockRepo.AssertExpectations(t)
}

func TestFindBySellerId(t *testing.T) {
	testProductSellerID = 1
	mockRepo.On("FindBySellerId", testProductSellerID).Return(&testProduct, nil)

	foundProduct, err := mockRepo.FindBySellerId(testProductSellerID)

	assert.NoError(t, err)
	assert.NotNil(t, foundProduct)

	assertProductEquality(t, &testProduct, foundProduct)

	mockRepo.AssertExpectations(t)
}
