package repositories

import (
	model "Dzaakk/simple-commerce/internal/product/models"
	"Dzaakk/simple-commerce/internal/shopping_cart/models"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ProductRepositoryImpl struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &ProductRepositoryImpl{
		DB: db,
	}
}

const (
	queryCreate                      = `INSERT INTO public.product (product_name, price, stock, category_id, seller_id, created, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	queryUpdate                      = `UPDATE public.product SET product_name=$1, price=$2, stock=$3, updated=NOW(), updated_by=$4 WHERE id=$5`
	queryFindByCategoryId            = `SELECT * FROM public.product WHERE category_id = $1`
	queryFindBySellerIdAndCategoryId = `SELECT * FROM public.product WHERE seller_id = $1 AND category_id = $2`
	queryFindBySellerId              = `SELECT * FROM public.product WHERE seller_id = $1`
	queryFindById                    = `SELECT * FROM public.product WHERE id = $1`
	queryGetPriceById                = `SELECT price FROM public.product WHERE id = $1`
	queryGetStockById                = `SELECT stock FROM public.product WHERE id = $1`
	queryFindByName                  = `SELECT * FROM public.product WHERE product_name like '%' || $1 || '%'`
	querySetStockById                = `UPDATE public.product SET stock = $1 WHERE id = $2`
	dbQueryTimeout                   = 3 * time.Second
	queryBase                        = "SELECT * FROM products WHERE 1=1"
)

func (repo *ProductRepositoryImpl) contextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, dbQueryTimeout)
}
func (repo *ProductRepositoryImpl) Create(ctx context.Context, data model.TProduct) (*model.TProduct, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, queryCreate, data.ProductName, data.Price, data.Stock, data.CategoryId, data.SellerId, data.Base.Created, data.Base.CreatedBy)
	if err != nil {
		return nil, response.ExecError("create product", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	data.Id = int(id)
	return &data, nil
}

func (repo *ProductRepositoryImpl) Update(ctx context.Context, data model.TProduct) (int64, error) {
	if data.Price < 0 {
		return 0, errors.New("invalid input parameter")
	}
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, queryUpdate, data.ProductName, data.Price, data.Stock, data.UpdatedBy, data.Id)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *ProductRepositoryImpl) FindById(ctx context.Context, id int) (*model.TProduct, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	rows, err := repo.DB.QueryContext(ctx, queryFindById, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	product, err := retrieveProduct(rows)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (repo *ProductRepositoryImpl) FindProductByFilters(ctx context.Context, categoryId, sellerId *int) ([]*model.TProduct, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	query := queryBase
	args := []interface{}{}

	if categoryId != nil {
		query += " AND category_id = ?"
		args = append(args, *categoryId)
	}
	if sellerId != nil {
		query += " AND seller_id = ?"
		args = append(args, *sellerId)
	}

	rows, err := repo.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	return scanProducts(rows)
}

func (repo *ProductRepositoryImpl) GetPriceById(ctx context.Context, id int) (float32, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	var balance float32
	err := repo.DB.QueryRowContext(ctx, queryGetPriceById, id).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (repo *ProductRepositoryImpl) GetStockById(ctx context.Context, id int) (int, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	var stock int
	err := repo.DB.QueryRowContext(ctx, queryGetPriceById, id).Scan(stock)
	if err != nil {
		return 0, err
	}

	return stock, nil
}

func (repo *ProductRepositoryImpl) SetStockById(ctx context.Context, id int, stock int) (int64, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, querySetStockById, stock, id)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *ProductRepositoryImpl) FindByName(ctx context.Context, name string) (*model.TProduct, error) {
	ctx, cancel := repo.contextWithTimeout(ctx)
	defer cancel()

	rows, err := repo.DB.QueryContext(ctx, queryFindByName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	product, err := retrieveProduct(rows)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (repo *ProductRepositoryImpl) UpdateStock(ctx context.Context, listData []*models.TCartItemDetail, name string) error {
	query, args := generateMultipleStockUpdateQuery(listData)
	_, err := repo.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepositoryImpl) UpdateStockWithTx(ctx context.Context, tx *sql.Tx, listItem []*models.TCartItemDetail) ([]*int, error) {
	listEmptyProductId, err := verifyStockAvailability(tx, listItem)
	if err != nil {
		return nil, err
	}
	query, args := generateMultipleStockUpdateQuery(listItem)
	result, err := tx.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error update product stock: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting rows affected: %w", err)
	}

	if int(rowsAffected) != len(listItem) {
		return nil, fmt.Errorf("expected to update %d products, but updated %d", len(listItem), rowsAffected)
	}
	return listEmptyProductId, nil
}
