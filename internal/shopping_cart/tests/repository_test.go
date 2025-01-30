package tests

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockRepo           = new(MockShoppingCartRepository)
	testShoppingCartID int
	testShoppingCart   = &model.TShoppingCart{
		Id:         1,
		CustomerId: 1,
		Status:     "A",
	}
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
