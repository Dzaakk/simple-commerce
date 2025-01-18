package repository

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"strconv"
	"time"
)

type ShoppingCartRepositoryImpl struct {
	DB *sql.DB
}

func NewShoppingCartRepository(db *sql.DB) ShoppingCartRepository {
	return &ShoppingCartRepositoryImpl{
		DB: db,
	}
}

const (
	queryCreateShoppingCart = `INSERT INTO public.shopping_cart (customer_id, status, created, created_by) VALUES ($1, $2, $3, $4) RETURNING id`
	queryFindById           = `SELECT id, customer_id, status FROM public.shopping_cart WHERE id=$1`
	queryFindByCartId       = `SELECT id, customer_id FROM public.shopping_cart WHERE customer_id=$1 AND status=$2`
	queryUpdateStatusById   = `UPDATE public.shopping_cart SET status=$1, updated_by=$2, updated=now() WHERE id=$3 RETURNING status`
	queryCheckStatus        = `SELECT status FROM public.shopping_cart WHERE id=$1 AND customer_id=$2`
	queryDeleteShoppingCart = `DELETE FROM public.shopping_cart WHERE id=$1`
	dbQueryTimeout          = 2 * time.Second
)

func (repo *ShoppingCartRepositoryImpl) contextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, dbQueryTimeout)
}

func (repo *ShoppingCartRepositoryImpl) Create(ctx context.Context, data model.TShoppingCart) (*model.TShoppingCart, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, queryCreateShoppingCart, data.CustomerId, data.Status, data.Base.Created, data.Base.CreatedBy)
	if err != nil {
		return nil, response.ExecError("create shopping cart", err)
	}

	id, _ := result.LastInsertId()
	newData := &model.TShoppingCart{
		Id:         int(id),
		CustomerId: data.CustomerId,
		Status:     data.Status,
	}

	return newData, nil
}

func (repo *ShoppingCartRepositoryImpl) FindById(ctx context.Context, id int) (*model.ShoppingCartRes, error) {
	var shoppingCart model.ShoppingCartRes
	err := repo.DB.QueryRow(queryFindById, id).Scan(&shoppingCart.Id, &shoppingCart.CustomerId, &shoppingCart.Status)
	if err != nil {
		return nil, err
	}

	return &shoppingCart, nil
}

func (repo *ShoppingCartRepositoryImpl) FindByCustomerIdAndStatus(ctx context.Context, customerId int, status string) (*model.TShoppingCart, error) {
	var data model.TShoppingCart
	err := repo.DB.QueryRow(queryFindByCartId, customerId, status).Scan(&data.Id, &data.CustomerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	data.Status = status
	return &data, nil
}

func (repo *ShoppingCartRepositoryImpl) UpdateStatusById(ctx context.Context, id int, status, customerId string) (*model.TShoppingCart, error) {
	statement, err := repo.DB.Prepare(queryUpdateStatusById)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var newStatus string
	err = statement.QueryRow(status, customerId, id).Scan(&newStatus)
	if err != nil {
		return nil, err
	}

	strCustomerId, _ := strconv.Atoi(customerId)
	shoppingCart := &model.TShoppingCart{
		Id:         id,
		CustomerId: strCustomerId,
		Status:     newStatus,
	}

	return shoppingCart, nil
}

func (repo *ShoppingCartRepositoryImpl) CheckStatus(ctx context.Context, id int, customerId int) (string, error) {
	var status string
	_ = repo.DB.QueryRow(queryCheckStatus, id, customerId).Scan(&status)

	return status, nil
}

func (repo *ShoppingCartRepositoryImpl) UpdateStatusByIdWithTx(ctx context.Context, tx *sql.Tx, cartId int, status string, customerid string) error {
	_, err := tx.Exec(queryUpdateStatusById, status, customerid, cartId)
	return err
}

func (repo *ShoppingCartRepositoryImpl) DeleteShoppingCart(ctx context.Context, cartId int) error {
	_, err := repo.DB.Exec(queryDeleteShoppingCart, cartId)
	return err
}
