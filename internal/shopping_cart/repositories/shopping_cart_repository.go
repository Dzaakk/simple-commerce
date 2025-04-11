package repositories

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	"context"
	"database/sql"
)

type ShoppingCartRepository interface {
	Create(ctx context.Context, data model.TShoppingCart) (*model.TShoppingCart, error)
	FindByStatusAndCustomerID(ctx context.Context, customerID int, status string) (*model.TShoppingCart, error)
	FindByCustomerID(ctx context.Context, customerID int) (*model.TShoppingCart, error)
	FindByCartID(ctx context.Context, cartID int) (*model.TShoppingCart, error)
	CheckStatus(ctx context.Context, cartID, customerID int) (string, error)
	UpdateStatusByCartID(ctx context.Context, cartID int, status, customerID string) (*model.TShoppingCart, error)
	UpdateStatusByCartIDWithTx(ctx context.Context, tx *sql.Tx, cartID int, status, customerID string) error
	DeleteShoppingCart(ctx context.Context, cartID int) error
}

type ShoppingCartItemRepository interface {
	Create(ctx context.Context, data model.TShoppingCartItem) (*model.TShoppingCartItem, error)
	Update(ctx context.Context, data model.TShoppingCartItem, customerID string) (*model.TShoppingCartItem, error)
	CountQuantityByProductIDAndCartID(ctx context.Context, productID, cartID int) (int, error)
	CountByCartID(ctx context.Context, cartID int) (int, error)
	Delete(ctx context.Context, productID, cartID int) error
	DeleteAll(ctx context.Context, cartID int) error
	DeleteAllWithTx(ctx context.Context, tx *sql.Tx, cartID int) error
	RetrieveCartItemsByCartID(ctx context.Context, cartID int) ([]*model.TCartItemDetail, error)
	RetrieveCartItemsByCartIDWithTx(ctx context.Context, tx *sql.Tx, cartID int) ([]*model.TCartItemDetail, error)
	SetEmptyQuantityWithTx(ctx context.Context, tx *sql.Tx, listProductID []*int) error
}
