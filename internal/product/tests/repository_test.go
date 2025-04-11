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
		ID:          1,
		ProductName: "Monitor",
		Price:       1700000,
		Stock:       10,
		CategoryID:  3,
		SellerID:    1,
	}
	testProduct2 = model.TProduct{
		ID:          2,
		ProductName: "Cooling Fan Ultra",
		Price:       650000,
		Stock:       20,
		CategoryID:  3,
		SellerID:    1,
	}
	testListProduct = []*model.TProduct{
		&testProduct, &testProduct2,
	}
	emptyListProduct      = []*model.TProduct{}
	testProductID         int
	testProductCategoryID int
	testProductName       string
	testProductSellerID   int
	testProductStock      int
	expectedError         error
)

func assertProductEquality(t *testing.T, expected, actual *model.TProduct) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.ProductName, actual.ProductName)
	assert.Equal(t, expected.Price, actual.Price)
	assert.Equal(t, expected.Stock, actual.Stock)
	assert.Equal(t, expected.CategoryID, actual.CategoryID)
	assert.Equal(t, expected.SellerID, actual.SellerID)
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
		testProduct.ID = 2
		expectedError = errors.New("failed to update product")
		mockRepo.On("Update", testProduct).Return(int64(0), expectedError)

		rowsAffected, err := mockRepo.Update(testProduct)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, int64(0), rowsAffected)
		mockRepo.AssertExpectations(t)
	})
}

func TestFindByCategoryID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testProductCategoryID = 3
		mockRepo.On("FindByCategoryID", testProductCategoryID).Return(testListProduct, nil)

		foundListProduct, err := mockRepo.FindByCategoryID(testProductCategoryID)

		assert.NoError(t, err)
		assert.NotNil(t, foundListProduct)
		assert.Equal(t, len(testListProduct), len(foundListProduct))

		for i, expectedProduct := range testListProduct {
			assertProductEquality(t, expectedProduct, foundListProduct[i])
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		testProductCategoryID = 99
		mockRepo.On("FindByCategoryID", testProductCategoryID).Return(emptyListProduct, nil)

		foundListProduct, err := mockRepo.FindByCategoryID(testProductCategoryID)

		assert.NoError(t, err)
		assert.Empty(t, foundListProduct)
		mockRepo.AssertExpectations(t)
	})

}

func TestFindByProductName(t *testing.T) {
	testProductName = "Monitor"
	mockRepo.On("FindByProductName", testProductName).Return(&testProduct, nil)

	foundProduct, err := mockRepo.FindByProductName(testProductName)

	assert.NoError(t, err)
	assert.NotNil(t, foundProduct)

	assertProductEquality(t, &testProduct, foundProduct)

	mockRepo.AssertExpectations(t)
}

func TestFindBySellerID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testProductSellerID = 1
		mockRepo.On("FindBySellerID", testProductSellerID).Return(testListProduct, nil)

		foundListProduct, err := mockRepo.FindBySellerID(testProductSellerID)

		assert.NoError(t, err)
		assert.NotNil(t, foundListProduct)
		assert.Equal(t, len(testListProduct), len(foundListProduct))

		for i, expectedProduct := range testListProduct {
			assertProductEquality(t, expectedProduct, foundListProduct[i])
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		testProductSellerID = 2
		mockRepo.On("FindBySellerID", testProductSellerID).Return(emptyListProduct, nil)

		foundListProduct, err := mockRepo.FindBySellerID(testProductSellerID)

		assert.NoError(t, err)
		assert.Empty(t, foundListProduct)
		mockRepo.AssertExpectations(t)
	})

}
func TestFindBySellerIDAndCategoryID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testProductSellerID = 1
		testProductCategoryID = 3
		mockRepo.On("FindBySellerIDAndCategoryID", testProductSellerID, testProductCategoryID).Return(testListProduct, nil)

		foundListProduct, err := mockRepo.FindBySellerIDAndCategoryID(testProductSellerID, testProductCategoryID)

		assert.NoError(t, err)
		assert.NotNil(t, foundListProduct)
		assert.Equal(t, len(testListProduct), len(foundListProduct))

		for i, expectedProduct := range testListProduct {
			assertProductEquality(t, expectedProduct, foundListProduct[i])
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		testProductSellerID = 2
		testProductCategoryID = 4
		mockRepo.On("FindBySellerIDAndCategoryID", testProductSellerID, testProductCategoryID).Return(emptyListProduct, nil)

		foundListProduct, err := mockRepo.FindBySellerIDAndCategoryID(testProductSellerID, testProductCategoryID)

		assert.NoError(t, err)
		assert.Empty(t, foundListProduct)
		mockRepo.AssertExpectations(t)
	})

}

func TestSetStockByProductID(t *testing.T) {
	testProductStock = 10
	t.Run("Success", func(t *testing.T) {
		testProductID = 1

		mockRepo.On("SetStockByProductID", testProductID, testProductStock).Return(int64(1), nil)

		rowsAffected, err := mockRepo.SetStockByProductID(testProductID, testProductStock)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), rowsAffected)
		mockRepo.AssertExpectations(t)
	})
	t.Run("Failed", func(t *testing.T) {
		testProductID = 2
		expectedError = errors.New("failed to update stock product")
		mockRepo.On("SetStockByProductID", testProductID, testProductStock).Return(int64(0), expectedError)

		rowsAffected, err := mockRepo.SetStockByProductID(testProductID, testProductStock)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, int64(0), rowsAffected)
		mockRepo.AssertExpectations(t)
	})
}
