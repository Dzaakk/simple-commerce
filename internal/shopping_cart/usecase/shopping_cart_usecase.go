package usecase

import (
	"Dzaakk/simple-commerce/internal/shopping_cart/model"
	"context"
)

type ShoppingCartUseCase interface {
	Add(ctx context.Context, data model.ShoppingCartReq) (*model.ShoppingCartItem, error)
	AddV2(ctx context.Context, data model.ShoppingCartReq) (*model.ShoppingCartItem, error)
	CheckStatus(ctx context.Context, cartID, customerID int) (string, error)
	CreateCart(ctx context.Context, customerID int) (*model.TShoppingCart, error)
	CreateCartItem(ctx context.Context, data model.ShoppingCartReq) (*model.TShoppingCartItem, error)
	GetListItem(ctx context.Context, customerID int) ([]*model.ListCartItemRes, error)
	DeleteShoppingList(ctx context.Context, data model.DeleteReq) error
}
