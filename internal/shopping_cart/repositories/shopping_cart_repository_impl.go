package repositories

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"errors"
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
	queryCreateShoppingCart        = `INSERT INTO public.shopping_cart (customer_id, status, created, created_by) VALUES ($1, $2, $3, $4) RETURNING id`
	queryFindById                  = `SELECT id, customer_id, status FROM public.shopping_cart WHERE id=$1`
	queryFindByCustomerIdAndStatus = `SELECT id, customer_id FROM public.shopping_cart WHERE customer_id=$1 AND status=$2`
	queryUpdateStatusById          = `UPDATE public.shopping_cart SET status=$1, updated_by=$2, updated=now() WHERE id=$3`
	queryCheckStatus               = `SELECT status FROM public.shopping_cart WHERE id=$1 AND customer_id=$2`
	queryDeleteShoppingCart        = `DELETE FROM public.shopping_cart WHERE id=$1`
	dbQueryTimeout                 = 2 * time.Second
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

func (repo *ShoppingCartRepositoryImpl) FindById(ctx context.Context, id int) (*model.TShoppingCart, error) {
	if id <= 0 {
		return nil, errors.New("invalid input parameter")
	}

	return repo.findShoppingCart(ctx, queryFindById, id)
}

func (repo *ShoppingCartRepositoryImpl) FindByCustomerIdAndStatus(ctx context.Context, customerId int, status string) (*model.TShoppingCart, error) {
	if customerId <= 0 || status == "" {
		return nil, errors.New("invalid input parameter")
	}
	return repo.findShoppingCart(ctx, queryFindByCustomerIdAndStatus, customerId, status)
}

func (repo *ShoppingCartRepositoryImpl) findShoppingCart(ctx context.Context, query string, args ...interface{}) (*model.TShoppingCart, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	row := repo.DB.QueryRowContext(ctx, query, args...)
	return scanCart(row)
}

func (repo *ShoppingCartRepositoryImpl) UpdateStatusById(ctx context.Context, id int, status, customerId string) (*model.TShoppingCart, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, queryUpdateStatusById, status, customerId, id)
	if err != nil {
		return nil, response.ExecError("update status by id", err)
	}

	intCustomerID, _ := strconv.Atoi(customerId)
	shoppingCart := &model.TShoppingCart{
		Id:         id,
		CustomerId: intCustomerID,
		Status:     status,
	}

	return shoppingCart, nil
}

func (repo *ShoppingCartRepositoryImpl) CheckStatus(ctx context.Context, id int, customerId int) (string, error) {
	if id <= 0 || customerId <= 0 {
		return "", errors.New("invalid input parameter")
	}

	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	var status string
	_ = repo.DB.QueryRowContext(ctx, queryCheckStatus, id, customerId).Scan(&status)

	return status, nil
}

func (repo *ShoppingCartRepositoryImpl) UpdateStatusByIdWithTx(ctx context.Context, tx *sql.Tx, cartId int, status string, customerid string) error {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	_, err := tx.ExecContext(ctx, queryUpdateStatusById, status, customerid, cartId)
	return err
}

func (repo *ShoppingCartRepositoryImpl) DeleteShoppingCart(ctx context.Context, cartId int) error {
	if cartId <= 0 {
		return errors.New("invalid input parameter")
	}
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, queryDeleteShoppingCart, cartId)
	if err != nil {
		return response.ExecError("delete shopping cart", err)
	}

	return nil
}
