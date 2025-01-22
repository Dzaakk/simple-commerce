package repositories

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
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
	queryCountItemByChartId     = `SELECT COUNT(*) FROM public.shopping_cart_item WHERE cart_id=$1`
	queryUpdateShoppingCartItem = `UPDATE public.shopping_cart_item SET quantity=$1, updated=now(), updated_by=$2 WHERE cart_id=$3 AND product_id=$4 RETURNING quantity`
	queryCreateShoppingCartItem = `INSERT INTO public.shopping_cart_item (cart_id, product_id, quantity, created, created_by) VALUES ($1, $2, $3, $4, $5)`
	queryCountProductQuantity   = `SELECT SUM(quantity) FROM public.shopping_cart_item WHERE product_id=$1 AND cart_id=$2`
	queryRetrieveCartItems      = "SELECT sci.product_id, p.product_name, p.price, sci.quantity FROM public.shopping_cart_item sci JOIN public.product p ON sci.product_id = p.id WHERE sci.cart_id=$1 ORDER BY p.product_name ASC"
	queryDeleteCartItems        = "DELETE FROM shopping_cart_item WHERE cart_id=$1 AND product_id=$2"
	dbQueryItemTimeout          = 2 * time.Second
)

func (repo *ShoppingCartItemRepositoryImpl) contextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, dbQueryItemTimeout)
}

func (repo *ShoppingCartItemRepositoryImpl) Create(ctx context.Context, data model.TShoppingCartItem) (*model.TShoppingCartItem, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, queryCreateShoppingCartItem, data.CartId, data.ProductId, data.Quantity, data.Base.Created, data.Base.CreatedBy)
	if err != nil {
		return nil, response.ExecError("create cart item", err)
	}

	return &data, nil
}

func (repo *ShoppingCartItemRepositoryImpl) SetEmptyQuantityWithTx(ctx context.Context, tx *sql.Tx, listProductId []*int) error {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	if len(listProductId) == 0 {
		return nil
	}
	var query strings.Builder

	query.WriteString("UPDATE public.shopping_cart_item SET quantity = 0 WHERE product_id IN (")

	for i := range listProductId {
		if i > 0 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("$%d", i+1))
	}
	query.WriteString(");")

	args := make([]interface{}, len(listProductId))
	for i, id := range listProductId {
		args[i] = *id
	}
	_, err := tx.ExecContext(ctx, query.String(), args...)
	if err != nil {
		return response.ExecError("set empty quantity with tx", err)
	}
	return nil
}

func (repo *ShoppingCartItemRepositoryImpl) DeleteAll(ctx context.Context, cartId int) error {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	if cartId <= 0 {
		return response.InvalidParameter()
	}

	_, err := repo.DB.ExecContext(ctx, queryDeleteAllCartItems, cartId)
	if err != nil {
		return response.ExecError("delete all", err)
	}

	return nil
}

func (repo *ShoppingCartItemRepositoryImpl) DeleteAllWithTx(ctx context.Context, tx *sql.Tx, cartId int) error {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	if cartId <= 0 {
		return response.InvalidParameter()
	}

	_, err := tx.ExecContext(ctx, queryDeleteAllCartItems, cartId)
	if err != nil {
		return response.ExecError("delete all with tx", err)
	}

	return nil
}

func (repo *ShoppingCartItemRepositoryImpl) CountByCartId(ctx context.Context, cartId int) (int, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	if cartId <= 0 {
		return 0, response.InvalidParameter()
	}

	var total int
	_ = repo.DB.QueryRowContext(ctx, queryCountItemByChartId, cartId).Scan(&total)
	return total, nil
}

func (repo *ShoppingCartItemRepositoryImpl) Update(ctx context.Context, data model.TShoppingCartItem, customerId string) (*model.TShoppingCartItem, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	custId, err := strconv.Atoi(customerId)
	if custId <= 0 || err != nil {
		return nil, response.InvalidParameter()
	}

	var updatedQuantity int
	_ = repo.DB.QueryRowContext(ctx, queryUpdateShoppingCartItem, data.Quantity, customerId, data.CartId, data.ProductId).Scan(&updatedQuantity)

	resData := &model.TShoppingCartItem{
		ProductId: data.ProductId,
		CartId:    data.CartId,
		Quantity:  updatedQuantity,
	}

	return resData, nil
}

func (repo *ShoppingCartItemRepositoryImpl) Delete(ctx context.Context, productId int, cartId int) error {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	if cartId <= 0 || productId <= 0 {
		return response.InvalidParameter()
	}

	_, err := repo.DB.ExecContext(ctx, queryDeleteCartItems, cartId, productId)
	if err != nil {
		return response.ExecError("delete cart item", err)
	}

	return nil
}

func (repo *ShoppingCartItemRepositoryImpl) RetrieveCartItemsByCartId(ctx context.Context, cartId int) ([]*model.TCartItemDetail, error) {

	rows, err := repo.DB.Query(queryRetrieveCartItems, cartId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []*model.TCartItemDetail
	for rows.Next() {
		var ci model.TCartItemDetail
		err := rows.Scan(&ci.ProductId, &ci.ProductName, &ci.Price, &ci.Quantity)
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

func (repo *ShoppingCartItemRepositoryImpl) RetrieveCartItemsByCartIdWithTx(ctx context.Context, tx *sql.Tx, cartId int) ([]*model.TCartItemDetail, error) {
	rows, err := tx.Query(queryRetrieveCartItems, cartId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []*model.TCartItemDetail
	for rows.Next() {
		var ci model.TCartItemDetail
		err := rows.Scan(&ci.ProductId, &ci.ProductName, &ci.Price, &ci.Quantity)
		if err != nil {
			return nil, err
		}
		if ci.Quantity == 0 {
			return nil, fmt.Errorf("product ID %d has a quantity of zero", ci.ProductId)
		}

		cartItems = append(cartItems, &ci)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (repo *ShoppingCartItemRepositoryImpl) CountQuantityByProductAndCartId(ctx context.Context, productId int, cartId int) (int, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	if productId <= 0 || cartId <= 0 {
		return 0, response.InvalidParameter()
	}

	var totalQuantity int
	_ = repo.DB.QueryRowContext(ctx, queryCountProductQuantity, productId, cartId).Scan(&totalQuantity)

	return totalQuantity, nil
}
