package tests

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockRepo         = new(MockShoppingCartRepository)
	testShoppingCart = &model.TShoppingCart{
		ID:         1,
		CustomerID: 1,
		Status:     "A",
	}
	testShoppingCartID, testCustomerID int
	testShoppingCartStatus             string
	ctx                                = context.Background()
)

func TestCreateShoppingCart(t *testing.T) {
	testShoppingCartID = 1

	mockRepo.On("Create", ctx, testShoppingCart).Return(&testShoppingCart, nil)

	createdShoppingCart, err := mockRepo.Create(ctx, testShoppingCart)

	assert.NoError(t, err)
	assert.NotNil(t, createdShoppingCart)

	assert.Equal(t, testShoppingCart.ID, createdShoppingCart.ID)
	assert.Equal(t, testShoppingCart.CustomerID, createdShoppingCart.CustomerID)
	assert.Equal(t, testShoppingCart.Status, createdShoppingCart.Status)

	mockRepo.AssertExpectations(t)
}
func TestFindByID(t *testing.T) {
	testShoppingCartID = 1

	mockRepo.On("FindByID", ctx, testShoppingCartID).Return(testShoppingCart, nil)
	foundShopingCart, err := mockRepo.FindByID(ctx, testShoppingCartID)

	assert.NoError(t, err)
	assert.NotNil(t, foundShopingCart)

	assert.Equal(t, testShoppingCart.ID, foundShopingCart.ID)
	assert.Equal(t, testShoppingCart.CustomerID, foundShopingCart.CustomerID)
	assert.Equal(t, testShoppingCart.Status, foundShopingCart.Status)

	mockRepo.AssertExpectations(t)
}

func TestFindByStatusAndCustomerID(t *testing.T) {
	testCustomerID = 1
	testShoppingCartStatus = "A"

	mockRepo.On("FindByStatusAndCustomerID", ctx, testShoppingCartStatus, testCustomerID).Return(testShoppingCart, nil)
	foundShopingCart, err := mockRepo.FindByStatusAndCustomerID(ctx, testShoppingCartStatus, testCustomerID)

	assert.NoError(t, err)
	assert.NotNil(t, foundShopingCart)

	assert.Equal(t, testShoppingCart.ID, foundShopingCart.ID)
	assert.Equal(t, testShoppingCart.CustomerID, foundShopingCart.CustomerID)
	assert.Equal(t, testShoppingCart.Status, foundShopingCart.Status)

	mockRepo.AssertExpectations(t)
}
