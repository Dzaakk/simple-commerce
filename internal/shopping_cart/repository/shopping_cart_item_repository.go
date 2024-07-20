package shopping_cart

import (
	model "Dzaakk/synapsis/internal/shopping_cart/models"
	"database/sql"
	"fmt"
)

type ShoppingCartItemRepository interface {
	Create(data model.TShoppingCartItem) (*model.TShoppingCartItem, error)
	Update(data model.TShoppingCartItem, customerId string) (*model.ShoppingCartItemRes, error)
	CountQuantityByProductAndCartId(productId, cartId int) (int, error)
	CountByCartId(cartId int) (int, error)
	Delete(productId, cartId int) error
	DeleteAll(cartId int) error
	RetrieveCartItemsByCartId(cartId int) ([]*model.TCartItemDetail, error)
}

type ShoppingCartItemRepositoryImpl struct {
	DB *sql.DB
}

func NewShoppingCartItemRepository(db *sql.DB) ShoppingCartItemRepository {
	return &ShoppingCartItemRepositoryImpl{
		DB: db,
	}
}

const queryDeleteAllCartItems = "DELETE FROM shopping_cart_item WHERE cart_id=$1"

func (repo *ShoppingCartItemRepositoryImpl) DeleteAll(cartId int) error {
	fmt.Println("CART ID = ", cartId)
	result, err := repo.DB.Exec(queryDeleteAllCartItems, cartId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("Number of rows deleted: %d\n", rowsAffected)
	return nil
}

const queryCountItemByChartId = `SELECT COUNT(*) FROM public.shopping_cart_item WHERE cart_id=$1`

func (repo *ShoppingCartItemRepositoryImpl) CountByCartId(cartId int) (int, error) {
	var total int
	err := repo.DB.QueryRow(queryCountItemByChartId, cartId).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

const queryUpdateShoppingCartItem = `UPDATE public.shopping_cart_item SET quantity=$1, updated=now(), updated_by=$2 WHERE cart_id=$3 AND product_id=$4 RETURNING quantity`

func (repo *ShoppingCartItemRepositoryImpl) Update(data model.TShoppingCartItem, customerId string) (*model.ShoppingCartItemRes, error) {
	statement, err := repo.DB.Prepare(queryUpdateShoppingCartItem)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var updatedQuantity int
	err = statement.QueryRow(data.Quantity, customerId, data.CartId, data.ProductId).Scan(&updatedQuantity)
	if err != nil {
		return nil, err
	}

	updatedCartItem := &model.ShoppingCartItemRes{
		CartId:    fmt.Sprintf("%d", data.CartId),
		ProductId: fmt.Sprintf("%d", data.ProductId),
		Quantity:  fmt.Sprintf("%d", updatedQuantity),
	}

	return updatedCartItem, nil
}

const queryDeleteCartItems = "DELETE FROM shopping_cart_item WHERE cart_id=$1 AND product_id=$2"

func (repo *ShoppingCartItemRepositoryImpl) Delete(productId int, cartId int) error {
	result, err := repo.DB.Exec(queryDeleteCartItems, cartId, productId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("Number of rows deleted: %d\n", rowsAffected)
	return nil
}

const queryRetrieveCartItems = "SELECT p.product_name, p.price, sci.quantity FROM public.shopping_cart_item sci JOIN public.product p ON sci.product_id = p.id WHERE sci.cart_id=$1 ORDER BY p.product_name ASC"

func (repo *ShoppingCartItemRepositoryImpl) RetrieveCartItemsByCartId(cartId int) ([]*model.TCartItemDetail, error) {
	rows, err := repo.DB.Query(queryRetrieveCartItems, cartId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []*model.TCartItemDetail
	for rows.Next() {
		var ci model.TCartItemDetail
		err := rows.Scan(&ci.ProductName, &ci.Price, &ci.Quantity)
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

const queryCountProductQuantity = `SELECT SUM(quantity) FROM public.shopping_cart_item WHERE product_id=$1 AND cart_id=$2`

func (repo *ShoppingCartItemRepositoryImpl) CountQuantityByProductAndCartId(productId int, cartId int) (int, error) {
	var totalQuantity int
	err := repo.DB.QueryRow(queryCountProductQuantity, productId, cartId).Scan(&totalQuantity)
	if err != nil {
		if totalQuantity == 0 {
			return 0, nil
		} else {
			return 0, err
		}
	}
	return totalQuantity, nil
}

const queryCreateShoppingCartItem = `INSERT INTO public.shopping_cart_item (cart_id, product_id, quantity, created, created_by) VALUES ($1, $2, $3, $4, $5)`

func (repo *ShoppingCartItemRepositoryImpl) Create(data model.TShoppingCartItem) (*model.TShoppingCartItem, error) {
	statement, err := repo.DB.Prepare(queryCreateShoppingCartItem)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	_, err = statement.Exec(data.CartId, data.ProductId, data.Quantity, data.Base.Created, data.Base.CreatedBy)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
