package repository

import (
	"Dzaakk/simple-commerce/internal/shopping_cart/model"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"strconv"
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
)

func (repo *ShoppingCartRepositoryImpl) Create(ctx context.Context, data model.TShoppingCart) (*model.TShoppingCart, error) {

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

	row := repo.DB.QueryRowContext(ctx, queryFindByCustomerId, customerID)

	return scanCart(row)
}

func (repo *ShoppingCartRepositoryImpl) FindByCartID(ctx context.Context, cartID int) (*model.TShoppingCart, error) {

	row := repo.DB.QueryRowContext(ctx, queryFindById, cartID)

	return scanCart(row)
}

func (repo *ShoppingCartRepositoryImpl) FindByStatusAndCustomerID(ctx context.Context, customerID int, status string) (*model.TShoppingCart, error) {

	row := repo.DB.QueryRowContext(ctx, queryFindByCustomerIdAndStatus, customerID, status)

	return scanCart(row)
}

func (repo *ShoppingCartRepositoryImpl) UpdateStatusByCartID(ctx context.Context, cartID int, status, customerID string) (*model.TShoppingCart, error) {

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

	var status string
	_ = repo.DB.QueryRowContext(ctx, queryCheckStatus, cartID, customerID).Scan(&status)

	return status, nil
}

func (repo *ShoppingCartRepositoryImpl) UpdateStatusByCartIDWithTx(ctx context.Context, tx *sql.Tx, cartID int, status string, customerID string) error {

	_, err := tx.ExecContext(ctx, queryUpdateStatusById, status, customerID, cartID)
	return err
}

func (repo *ShoppingCartRepositoryImpl) DeleteShoppingCart(ctx context.Context, cartID int) error {

	_, err := repo.DB.ExecContext(ctx, queryDeleteShoppingCart, cartID)
	if err != nil {
		return response.ExecError("delete shopping cart", err)
	}

	return nil
}
