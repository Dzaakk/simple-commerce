package usecase

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	"context"
)

type ShoppingCartUseCase interface {
	Add(ctx context.Context, data model.ShoppingCartReq) (*model.ShoppingCartItem, error)
	CheckStatus(ctx context.Context, id, customerId int) (string, error)
	CreateCart(ctx context.Context, customerId int) (*model.TShoppingCart, error)
	CreateCartItem(ctx context.Context, data model.ShoppingCartReq) (*model.TShoppingCartItem, error)
	GetListItem(ctx context.Context, customerId int) ([]*model.ListCartItemRes, error)
	DeleteShoppingList(ctx context.Context, data model.DeleteReq) error
}
