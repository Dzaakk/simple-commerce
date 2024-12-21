package tests

import (
	model "Dzaakk/simple-commerce/internal/product/models"
	"errors"
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
	testProductStock    int
	expectedError       error
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
func TestUpdateProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testProduct.ProductName = "Laptop"

		mockRepo.On("Update", testProduct).Return(int64(1), nil)

		rowsAffected, err := mockRepo.Update(testProduct)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), rowsAffected)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		testProduct.Id = 2
		expectedError = errors.New("failed to update product")
		mockRepo.On("Update", testProduct).Return(int64(0), expectedError)

		rowsAffected, err := mockRepo.Update(testProduct)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, int64(0), rowsAffected)
		mockRepo.AssertExpectations(t)
	})
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
func TestSetStockById(t *testing.T) {
	testProductStock = 10
	t.Run("Success", func(t *testing.T) {
		testProductID = 1

		mockRepo.On("SetStockById", testProductID, testProductStock).Return(int64(1), nil)

		rowsAffected, err := mockRepo.SetStockById(testProductID, testProductStock)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), rowsAffected)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Failed", func(t *testing.T) {
		testProductID = 2
		expectedError = errors.New("failed to update stock product")
		mockRepo.On("SetStockById", testProductID, testProductStock).Return(int64(0), expectedError)

		rowsAffected, err := mockRepo.SetStockById(testProductID, testProductStock)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, int64(0), rowsAffected)
		mockRepo.AssertExpectations(t)
	})
}
