package repository

import (
	"Dzaakk/simple-commerce/internal/shopping_cart/model"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type ShoppingCartItemRepositoryImpl struct {
	DB *sql.DB
}

func NewShoppingCartItemRepository(db *sql.DB) ShoppingCartItemRepository {
	return &ShoppingCartItemRepositoryImpl{
		DB: db,
	}
}

const (
	queryDeleteAllCartItems     = "DELETE FROM shopping_cart_item WHERE cart_id=$1"
	queryCountItemByCartId      = `SELECT COUNT(*) FROM public.shopping_cart_item WHERE cart_id=$1`
	queryUpdateShoppingCartItem = `UPDATE public.shopping_cart_item SET quantity=$1, updated=now(), updated_by=$2 WHERE cart_id=$3 AND product_id=$4 RETURNING quantity`
	queryCreateShoppingCartItem = `INSERT INTO public.shopping_cart_item (cart_id, product_id, quantity, created, created_by) VALUES ($1, $2, $3, $4, $5)`
	queryCountProductQuantity   = `SELECT SUM(quantity) FROM public.shopping_cart_item WHERE product_id=$1 AND cart_id=$2`
	queryRetrieveCartItems      = "SELECT sci.product_id, p.product_name, p.price, sci.quantity FROM public.shopping_cart_item sci JOIN public.product p ON sci.product_id = p.id WHERE sci.cart_id=$1 ORDER BY p.product_name ASC"
	queryDeleteCartItems        = "DELETE FROM shopping_cart_item WHERE cart_id=$1 AND product_id=$2"
)

func (repo *ShoppingCartItemRepositoryImpl) Create(ctx context.Context, data model.TShoppingCartItem) (*model.TShoppingCartItem, error) {

	_, err := repo.DB.ExecContext(ctx, queryCreateShoppingCartItem, data.CartID, data.ProductID, data.Quantity, data.Base.Created, data.Base.CreatedBy)
	if err != nil {
		return nil, response.ExecError("create cart item", err)
	}

	return &data, nil
}

func (repo *ShoppingCartItemRepositoryImpl) SetEmptyQuantityWithTx(ctx context.Context, tx *sql.Tx, listProductID []*int) error {

	if len(listProductID) == 0 {
		return nil
	}

	var query strings.Builder

	query.WriteString("UPDATE public.shopping_cart_item SET quantity = 0 WHERE product_id IN (")

	for i := range listProductID {
		if i > 0 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("$%d", i+1))
	}
	query.WriteString(");")

	args := make([]interface{}, len(listProductID))
	for i, id := range listProductID {
		args[i] = *id
	}
	_, err := tx.ExecContext(ctx, query.String(), args...)
	if err != nil {
		return response.ExecError("set empty quantity with tx", err)
	}
	return nil
}

func (repo *ShoppingCartItemRepositoryImpl) DeleteAll(ctx context.Context, cartID int) error {

	_, err := repo.DB.ExecContext(ctx, queryDeleteAllCartItems, cartID)
	if err != nil {
		return response.ExecError("delete all", err)
	}

	return nil
}

func (repo *ShoppingCartItemRepositoryImpl) DeleteAllWithTx(ctx context.Context, tx *sql.Tx, cartID int) error {

	_, err := tx.ExecContext(ctx, queryDeleteAllCartItems, cartID)
	if err != nil {
		return response.ExecError("delete all with tx", err)
	}

	return nil
}

func (repo *ShoppingCartItemRepositoryImpl) CountByCartID(ctx context.Context, cartID int) (int, error) {

	var total int
	err := repo.DB.QueryRowContext(ctx, queryCountItemByCartId, cartID).Scan(&total)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, response.ExecError("count by cart id", err)
	}
	return total, nil
}

func (repo *ShoppingCartItemRepositoryImpl) Update(ctx context.Context, data model.TShoppingCartItem, customerID string) (*model.TShoppingCartItem, error) {

	var updatedQuantity int
	err := repo.DB.QueryRowContext(ctx, queryUpdateShoppingCartItem, data.Quantity, customerID, data.CartID, data.ProductID).Scan(&updatedQuantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.Error("cart item", err)
		}
		return nil, response.ExecError("update cart item", err)
	}
	resData := &model.TShoppingCartItem{
		ProductID: data.ProductID,
		CartID:    data.CartID,
		Quantity:  updatedQuantity,
	}

	return resData, nil
}

func (repo *ShoppingCartItemRepositoryImpl) Delete(ctx context.Context, productID int, cartID int) error {

	_, err := repo.DB.ExecContext(ctx, queryDeleteCartItems, cartID, productID)
	if err != nil {
		return response.ExecError("delete cart item", err)
	}

	return nil
}

func (repo *ShoppingCartItemRepositoryImpl) RetrieveCartItemsByCartID(ctx context.Context, cartID int) ([]*model.TCartItemDetail, error) {

	rows, err := repo.DB.Query(queryRetrieveCartItems, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []*model.TCartItemDetail
	for rows.Next() {
		var ci model.TCartItemDetail
		err := rows.Scan(&ci.ProductID, &ci.ProductName, &ci.Price, &ci.Quantity)
		if err != nil {
			return nil, err
		}
		cartItems = append(cartItems, &ci)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (repo *ShoppingCartItemRepositoryImpl) RetrieveCartItemsByCartIDWithTx(ctx context.Context, tx *sql.Tx, cartID int) ([]*model.TCartItemDetail, error) {
	rows, err := tx.Query(queryRetrieveCartItems, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []*model.TCartItemDetail
	for rows.Next() {
		var ci model.TCartItemDetail
		err := rows.Scan(&ci.ProductID, &ci.ProductName, &ci.Price, &ci.Quantity)
		if err != nil {
			return nil, err
		}
		if ci.Quantity == 0 {
			return nil, fmt.Errorf("product ID %d has a quantity of zero", ci.ProductID)
		}

		cartItems = append(cartItems, &ci)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (repo *ShoppingCartItemRepositoryImpl) CountQuantityByProductIDAndCartID(ctx context.Context, productID int, cartID int) (int, error) {

	var totalQuantity int
	_ = repo.DB.QueryRowContext(ctx, queryCountProductQuantity, productID, cartID).Scan(&totalQuantity)

	return totalQuantity, nil
}
