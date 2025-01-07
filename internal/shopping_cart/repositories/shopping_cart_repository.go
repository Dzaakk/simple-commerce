package repository

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	"database/sql"
)

type ShoppingCartRepository interface {
	Create(data model.TShoppingCart) (*model.TShoppingCart, error)
	FindByCustomerIdAndStatus(customerId int, status string) (*model.TShoppingCart, error)
	FindById(id int) (*model.ShoppingCartRes, error)
	CheckStatus(id, customerId int) (string, error)
	UpdateStatusById(id int, status, customerid string) (*model.TShoppingCart, error)
	UpdateStatusByIdWithTx(tx *sql.Tx, cartId int, status, customerid string) error
	DeleteShoppingCart(cartId int) error
}

type ShoppingCartItemRepository interface {
	Create(data model.TShoppingCartItem) (*model.TShoppingCartItem, error)
	Update(data model.TShoppingCartItem, customerId string) (*model.ShoppingCartItemRes, error)
	CountQuantityByProductAndCartId(productId, cartId int) (int, error)
	CountByCartId(cartId int) (int, error)
	Delete(productId, cartId int) error
	DeleteAll(cartId int) error
	DeleteAllWithTx(tx *sql.Tx, cartId int) error
	RetrieveCartItemsByCartId(cartId int) ([]*model.TCartItemDetail, error)
	RetrieveCartItemsByCartIdWithTx(tx *sql.Tx, cartId int) ([]*model.TCartItemDetail, error)
	SetQuantityWithTx(tx *sql.Tx, listProductId []*int) error
}
