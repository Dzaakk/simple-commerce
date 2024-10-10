package usecase

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
)

type ShoppingCartUseCase interface {
	Add(data model.ShoppingCartReq) (*model.ShoppingCartItem, error)
	CheckStatus(id, customerId int) (*string, error)
	CreateCart(customerId int) (*model.TShoppingCart, error)
	CreateCartItem(data model.ShoppingCartReq) (*model.TShoppingCartItem, error)
	GetListItem(customerId int) ([]*model.ListCartItemRes, error)
	DeleteShoppingList(data model.DeleteReq) error
}
