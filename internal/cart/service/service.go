package service

import (
	"Dzaakk/simple-commerce/internal/cart/dto"
	cartModel "Dzaakk/simple-commerce/internal/cart/model"
	catalogDto "Dzaakk/simple-commerce/internal/catalog/dto"
	catalogModel "Dzaakk/simple-commerce/internal/catalog/model"
	"context"
)

type CartService interface {
	GetCartItems(ctx context.Context, customerID string) (*dto.CartRes, error)
	AddItem(ctx context.Context, customerID string, productID string, quantity int) (*dto.CartRes, error)
	UpdateItem(ctx context.Context, customerID string, productID string, quantity int) (*dto.CartRes, error)
	DeleteItem(ctx context.Context, customerID string, productID string) error
	ClearItems(ctx context.Context, customerID string) error
}

type CartRepository interface {
	GetCartByCustomerID(ctx context.Context, customerID string) (*cartModel.Cart, error)
	GetOrCreateCart(ctx context.Context, customerID string) (*cartModel.Cart, error)
}

type CartItemRepository interface {
	GetCartItems(ctx context.Context, cartID string) ([]*cartModel.CartItem, error)
	UpsertItem(ctx context.Context, cartID string, productID string, quantity int, priceSnapshot float64) error
	DeleteItem(ctx context.Context, cartID string, productID string) error
	ClearItems(ctx context.Context, cartID string) error
}

type ProductService interface {
	FindByID(ctx context.Context, productID string) (*catalogDto.ProductRes, error)
}

type InventoryService interface {
	FindByProductID(ctx context.Context, productID string) (*catalogModel.Inventory, error)
}
