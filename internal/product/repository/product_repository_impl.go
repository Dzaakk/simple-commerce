package repository

import (
	"Dzaakk/simple-commerce/internal/product/model"
	cartModel "Dzaakk/simple-commerce/internal/shopping_cart/model"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"fmt"
	"strconv"
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
	queryFindByCategoryID            = `SELECT * FROM public.product WHERE category_id = $1`
	queryFindBySellerIDAndCategoryID = `SELECT * FROM public.product WHERE seller_id = $1 AND category_id = $2`
	queryFindBySellerID              = `SELECT * FROM public.product WHERE seller_id = $1`
	queryFindByProductID             = `SELECT * FROM public.product WHERE id = $1`
	queryGetPriceByProductID         = `SELECT price FROM public.product WHERE id = $1`
	queryGetStockByProductID         = `SELECT stock FROM public.product WHERE id = $1`
	queryFindByName                  = `SELECT * FROM public.product WHERE product_name like '%' || $1 || '%'`
	querySetStockByProductID         = `UPDATE public.product SET stock = $1 WHERE id = $2`
	queryBase                        = "SELECT * FROM products WHERE 1=1"
)

func (repo *ProductRepositoryImpl) Create(ctx context.Context, data model.TProduct) (*model.TProduct, error) {

	result, err := repo.DB.ExecContext(ctx, queryCreate, data.ProductName, data.Price, data.Stock, data.CategoryID, data.SellerID, data.Base.Created, data.Base.CreatedBy)
	if err != nil {
		return nil, response.ExecError("create product", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	data.ID = int(id)
	return &data, nil
}

func (repo *ProductRepositoryImpl) Update(ctx context.Context, data model.TProduct) (int64, error) {

	result, err := repo.DB.ExecContext(ctx, queryUpdate, data.ProductName, data.Price, data.Stock, data.UpdatedBy, data.ID)
	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *ProductRepositoryImpl) FindByFilters(ctx context.Context, params model.ProductFilter) ([]*model.TProduct, error) {
	baseQuery := `
		SELECT id, product_name, price, stock, category_id, seller_id
		FROM t_product
		WHERE 1=1
	`
	args := []interface{}{}
	i := 1

	if params.ProductName != "" {
		baseQuery += fmt.Sprintf(" AND product_name ILIKE $%d", i)
		args = append(args, "%"+params.ProductName+"%")
		i++
	}

	if params.CategoryID != "" {
		categoryID, err := strconv.Atoi(params.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %w", err)
		}
		baseQuery += fmt.Sprintf(" AND category_id = $%d", i)
		args = append(args, categoryID)
		i++
	}

	if params.SellerID != "" {
		sellerID, err := strconv.Atoi(params.SellerID)
		if err != nil {
			return nil, fmt.Errorf("invalid seller ID: %w", err)
		}
		baseQuery += fmt.Sprintf(" AND seller_id = $%d", i)
		args = append(args, sellerID)
		i++
	}

	if params.LowPrice != "" {
		lowPrice, err := strconv.ParseFloat(params.LowPrice, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid low price: %w", err)
		}
		baseQuery += fmt.Sprintf(" AND price >= $%d", i)
		args = append(args, lowPrice)
		i++
	}

	if params.HighPrice != "" {
		highPrice, err := strconv.ParseFloat(params.HighPrice, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid high price: %w", err)
		}
		baseQuery += fmt.Sprintf(" AND price <= $%d", i)
		args = append(args, highPrice)
		i++
	}

	baseQuery += " ORDER BY id DESC"

	if params.Limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", i)
		args = append(args, params.Limit)
		i++
	}

	if params.Offset >= 0 {
		baseQuery += fmt.Sprintf(" OFFSET $%d", i)
		args = append(args, params.Offset)
	}

	rows, err := repo.DB.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	return scanProducts(rows)
}

func (repo *ProductRepositoryImpl) FindByID(ctx context.Context, productID int) (*model.TProduct, error) {

	rows, err := repo.DB.QueryContext(ctx, queryFindByProductID, productID)
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

func (repo *ProductRepositoryImpl) FindByProductName(ctx context.Context, productName string) (*model.TProduct, error) {

	rows, err := repo.DB.QueryContext(ctx, queryFindByName, productName)
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

func (repo *ProductRepositoryImpl) UpdateStock(ctx context.Context, listData []*cartModel.TCartItemDetail, name string) error {

	query, args := generateMultipleStockUpdateQuery(listData)
	_, err := repo.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepositoryImpl) UpdateStockWithTx(ctx context.Context, tx *sql.Tx, listItem []*cartModel.TCartItemDetail) ([]*int, error) {
	listEmptyProductID, err := verifyStockAvailability(tx, listItem)
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
	return listEmptyProductID, nil
}

// func (repo *ProductRepositoryImpl) FindProductByFilters(ctx context.Context, categoryID, sellerID *int) ([]*model.TProduct, error) {

// 	query := queryBase
// 	args := []interface{}{}

// 	if categoryID != nil {
// 		query += " AND category_id = ?"
// 		args = append(args, *categoryID)
// 	}
// 	if sellerID != nil {
// 		query += " AND seller_id = ?"
// 		args = append(args, *sellerID)
// 	}

// 	rows, err := repo.DB.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute query: %w", err)
// 	}
// 	defer rows.Close()

// 	return scanProducts(rows)
// }

// func (repo *ProductRepositoryImpl) GetPriceByID(ctx context.Context, productID int) (float32, error) {

// 	var balance float32
// 	err := repo.DB.QueryRowContext(ctx, queryGetPriceByProductID, productID).Scan(&balance)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return balance, nil
// }

// func (repo *ProductRepositoryImpl) GetStockByID(ctx context.Context, productID int) (int, error) {

// 	var stock int
// 	err := repo.DB.QueryRowContext(ctx, queryGetPriceByProductID, productID).Scan(stock)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return stock, nil
// }

// func (repo *ProductRepositoryImpl) SetStockByID(ctx context.Context, productID int, stock int) (int64, error) {

// 	result, err := repo.DB.ExecContext(ctx, querySetStockByProductID, stock, productID)
// 	if err != nil {
// 		return 0, err
// 	}
// 	rowsAffected, _ := result.RowsAffected()
// 	return rowsAffected, nil
// }
