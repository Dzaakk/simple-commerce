package shopping_cart

import (
	model "Dzaakk/synapsis/internal/shopping_cart/models"
	"database/sql"
	"fmt"
	"strconv"
)

type ShoppingCartRepository interface {
	Create(data model.TShoppingCart) (*model.ShoppingCartRes, error)
	FindByCustomerIdAndStatus(customerId int, status string) (*model.ShoppingCartRes, error)
	FindById(id int) (*model.ShoppingCartRes, error)
	CheckStatus(id, customerId int) (*string, error)
	UpdateStatusById(id int, status string) (*model.ShoppingCartRes, error)
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

func (repo *ShoppingCartRepositoryImpl) Create(data model.TShoppingCart) (*model.ShoppingCartRes, error) {
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

	newShoppingCart := &model.ShoppingCartRes{
		Id:         fmt.Sprintf("%d", id),
		CustomerId: fmt.Sprintf("%d", data.CustomerId),
		Status:     data.Status,
	}

	return newShoppingCart, nil
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

func (repo *ShoppingCartRepositoryImpl) FindByCustomerIdAndStatus(customerId int, status string) (*model.ShoppingCartRes, error) {
	var data model.TShoppingCart
	err := repo.DB.QueryRow(queryFindByCartId, customerId, status).Scan(&data.Id, &data.CustomerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	shoppingCart := &model.ShoppingCartRes{
		Id:         fmt.Sprintf("%d", data.Id),
		CustomerId: fmt.Sprintf("%d", data.CustomerId),
		Status:     status,
	}
	return shoppingCart, nil
}

const queryUpdateStatusById = `UPDATE public.shopping_cart SET status=$1 updated_by=$2, updated=now() WHERE id=$3 RETURNING status, customer_id`

func (repo *ShoppingCartRepositoryImpl) UpdateStatusById(id int, status string) (*model.ShoppingCartRes, error) {
	statement, err := repo.DB.Prepare(queryUpdateStatusById)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var newStatus string
	var customerId int
	idString := strconv.Itoa(id)
	err = statement.QueryRow(status, idString, id).Scan(&newStatus, &customerId)
	if err != nil {
		return nil, err
	}

	shoppingCart := &model.ShoppingCartRes{
		Id:         fmt.Sprintf("%d", id),
		CustomerId: fmt.Sprintf("%d", customerId),
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
