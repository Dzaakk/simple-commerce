package tests

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockRepo         = new(MockShoppingCartRepository)
	testShoppingCart = &model.TShoppingCart{
		Id:         1,
		CustomerId: 1,
		Status:     "A",
	}
	testShoppingCartID, testCustomerID int
	testShoppingCartStatus             string
)

func TestFindByID(t *testing.T) {
	testShoppingCartID = 1

	mockRepo.On("FindByID", testShoppingCartID).Return(testShoppingCart, nil)
	foundShopingCart, err := mockRepo.FindByID(testShoppingCartID)

	assert.NoError(t, err)
	assert.NotNil(t, foundShopingCart)

	assert.Equal(t, testShoppingCart.Id, foundShopingCart.Id)
	assert.Equal(t, testShoppingCart.CustomerId, foundShopingCart.CustomerId)
	assert.Equal(t, testShoppingCart.Status, foundShopingCart.Status)

	mockRepo.AssertExpectations(t)
}

func TestFindByStatusAndCustomerID(t *testing.T) {
	testCustomerID = 1
	testShoppingCartStatus = "A"

	mockRepo.On("FindByStatusAndCustomerID", testShoppingCartStatus, testCustomerID).Return(testShoppingCart, nil)
	foundShopingCart, err := mockRepo.FindByStatusAndCustomerID(testShoppingCartStatus, testCustomerID)

	assert.NoError(t, err)
	assert.NotNil(t, foundShopingCart)

	assert.Equal(t, testShoppingCart.Id, foundShopingCart.Id)
	assert.Equal(t, testShoppingCart.CustomerId, foundShopingCart.CustomerId)
	assert.Equal(t, testShoppingCart.Status, foundShopingCart.Status)

	mockRepo.AssertExpectations(t)
}
