package shopping_cart

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	"database/sql"
	"fmt"
	"strconv"
)

type ShoppingCartRepository interface {
	Create(data model.TShoppingCart) (*model.TShoppingCart, error)
	FindByCustomerIdAndStatus(customerId int, status string) (*model.TShoppingCart, error)
	FindById(id int) (*model.ShoppingCartRes, error)
	CheckStatus(id, customerId int) (*string, error)
	UpdateStatusById(id int, status, customerid string) (*model.TShoppingCart, error)
}

type ShoppingCartRepositoryImpl struct {
	DB *sql.DB
}

func NewShoppingCartRepository(db *sql.DB) ShoppingCartRepository {
	return &ShoppingCartRepositoryImpl{
		DB: db,
	}
}

const queryCreateShoppingCart = `INSERT INTO public.shopping_cart (customer_id, status, created, created_by) VALUES ($1, $2, $3, $4) RETURNING id`

func (repo *ShoppingCartRepositoryImpl) Create(data model.TShoppingCart) (*model.TShoppingCart, error) {
	statement, err := repo.DB.Prepare(queryCreateShoppingCart)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var id int

	err = statement.QueryRow(data.CustomerId, data.Status, data.Base.Created, data.Base.CreatedBy).Scan(&id)
	if err != nil {
		return nil, err
	}

	newData := &model.TShoppingCart{
		Id:         id,
		CustomerId: data.CustomerId,
		Status:     data.Status,
	}

	return newData, nil
}

const queryFindById = `SELECT id, customer_id, status FROM public.shopping_cart WHERE id=$1`

func (repo *ShoppingCartRepositoryImpl) FindById(id int) (*model.ShoppingCartRes, error) {
	var shoppingCart model.ShoppingCartRes
	err := repo.DB.QueryRow(queryFindById, id).Scan(&shoppingCart.Id, &shoppingCart.CustomerId, &shoppingCart.Status)
	if err != nil {
		return nil, err
	}

	return &shoppingCart, nil
}

const queryFindByCartId = `SELECT id, customer_id FROM public.shopping_cart WHERE customer_id=$1 AND status=$2`

func (repo *ShoppingCartRepositoryImpl) FindByCustomerIdAndStatus(customerId int, status string) (*model.TShoppingCart, error) {
	var data model.TShoppingCart
	err := repo.DB.QueryRow(queryFindByCartId, customerId, status).Scan(&data.Id, &data.CustomerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	data.Status = status
	fmt.Println("DATA = ", data)
	return &data, nil
}

const queryUpdateStatusById = `UPDATE public.shopping_cart SET status=$1, updated_by=$2, updated=now() WHERE id=$3 RETURNING status`

func (repo *ShoppingCartRepositoryImpl) UpdateStatusById(id int, status, customerId string) (*model.TShoppingCart, error) {
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

const queryCheckStatus = `SELECT status FROM public.shopping_cart WHERE id=$1 AND customer_id=$2`

func (repo *ShoppingCartRepositoryImpl) CheckStatus(id int, customerId int) (*string, error) {
	var status string
	_ = repo.DB.QueryRow(queryCheckStatus, id, customerId).Scan(&status)

	return &status, nil
}
