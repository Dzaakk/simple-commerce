package repositories

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	"context"
	"database/sql"
)

type ShoppingCartRepository interface {
	Create(ctx context.Context, data model.TShoppingCart) (*model.TShoppingCart, error)
	FindByCustomerIdAndStatus(ctx context.Context, customerId int, status string) (*model.TShoppingCart, error)
	FindById(ctx context.Context, id int) (*model.TShoppingCart, error)
	CheckStatus(ctx context.Context, id, customerId int) (string, error)
	UpdateStatusById(ctx context.Context, id int, status, customerid string) (*model.TShoppingCart, error)
	UpdateStatusByIdWithTx(ctx context.Context, tx *sql.Tx, cartId int, status, customerid string) error
	DeleteShoppingCart(ctx context.Context, cartId int) error
}

type ShoppingCartItemRepository interface {
	Create(ctx context.Context, data model.TShoppingCartItem) (*model.TShoppingCartItem, error)
	Update(ctx context.Context, data model.TShoppingCartItem, customerId string) (*model.TShoppingCartItem, error)
	CountQuantityByProductAndCartId(ctx context.Context, productId, cartId int) (int, error)
	CountByCartId(ctx context.Context, cartId int) (int, error)
	Delete(ctx context.Context, productId, cartId int) error
	DeleteAll(ctx context.Context, cartId int) error
	DeleteAllWithTx(ctx context.Context, tx *sql.Tx, cartId int) error
	RetrieveCartItemsByCartId(ctx context.Context, cartId int) ([]*model.TCartItemDetail, error)
	RetrieveCartItemsByCartIdWithTx(ctx context.Context, tx *sql.Tx, cartId int) ([]*model.TCartItemDetail, error)
	SetEmptyQuantityWithTx(ctx context.Context, tx *sql.Tx, listProductId []*int) error
}
