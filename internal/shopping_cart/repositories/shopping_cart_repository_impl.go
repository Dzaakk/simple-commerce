package repositories

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
	queryCreateShoppingCart        = `INSERT INTO public.shopping_cart (customer_id, status, created, created_by) VALUES ($1, $2, $3, $4) RETURNING id`
	queryFindById                  = `SELECT id, customer_id, status FROM public.shopping_cart WHERE id=$1`
	queryFindByCustomerId          = `SELECT id, customer_id, status FROM public.shopping_cart WHERE customer_id=$1`
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

	result, err := repo.DB.ExecContext(ctx, queryCreateShoppingCart, data.CustomerID, data.Status, data.Base.Created, data.Base.CreatedBy)
	if err != nil {
		return nil, response.ExecError("create shopping cart", err)
	}

	cartID, _ := result.LastInsertId()
	newData := &model.TShoppingCart{
		ID:         int(cartID),
		CustomerID: data.CustomerID,
		Status:     data.Status,
	}

	return newData, nil
}
func (repo *ShoppingCartRepositoryImpl) FindByCustomerID(ctx context.Context, customerID int) (*model.TShoppingCart, error) {
	if customerID <= 0 {
		return nil, response.InvalidParameter()
	}

	return repo.findShoppingCart(ctx, queryFindByCustomerId, customerID)
}

func (repo *ShoppingCartRepositoryImpl) FindByCartID(ctx context.Context, cartID int) (*model.TShoppingCart, error) {
	if cartID <= 0 {
		return nil, response.InvalidParameter()
	}

	return repo.findShoppingCart(ctx, queryFindById, cartID)
}

func (repo *ShoppingCartRepositoryImpl) FindByStatusAndCustomerID(ctx context.Context, customerID int, status string) (*model.TShoppingCart, error) {
	if customerID <= 0 || status == "" {
		return nil, response.InvalidParameter()
	}
	return repo.findShoppingCart(ctx, queryFindByCustomerIdAndStatus, customerID, status)
}

func (repo *ShoppingCartRepositoryImpl) findShoppingCart(ctx context.Context, query string, args ...interface{}) (*model.TShoppingCart, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	row := repo.DB.QueryRowContext(ctx, query, args...)
	return scanCart(row)
}

func (repo *ShoppingCartRepositoryImpl) UpdateStatusByCartID(ctx context.Context, cartID int, status, customerID string) (*model.TShoppingCart, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, queryUpdateStatusById, status, customerID, cartID)
	if err != nil {
		return nil, response.ExecError("update status by id", err)
	}

	intCustomerID, _ := strconv.Atoi(customerID)
	shoppingCart := &model.TShoppingCart{
		ID:         cartID,
		CustomerID: intCustomerID,
		Status:     status,
	}

	return shoppingCart, nil
}

func (repo *ShoppingCartRepositoryImpl) CheckStatus(ctx context.Context, cartID int, customerID int) (string, error) {
	if cartID <= 0 || customerID <= 0 {
		return "", response.InvalidParameter()
	}

	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	var status string
	_ = repo.DB.QueryRowContext(ctx, queryCheckStatus, cartID, customerID).Scan(&status)

	return status, nil
}

func (repo *ShoppingCartRepositoryImpl) UpdateStatusByCartIDWithTx(ctx context.Context, tx *sql.Tx, cartID int, status string, customerID string) error {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	_, err := tx.ExecContext(ctx, queryUpdateStatusById, status, customerID, cartID)
	return err
}

func (repo *ShoppingCartRepositoryImpl) DeleteShoppingCart(ctx context.Context, cartID int) error {
	if cartID <= 0 {
		return response.InvalidParameter()
	}
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	_, err := repo.DB.ExecContext(ctx, queryDeleteShoppingCart, cartID)
	if err != nil {
		return response.ExecError("delete shopping cart", err)
	}

	return nil
}
