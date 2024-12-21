package repositories

import (
	model "Dzaakk/simple-commerce/internal/product/models"
	"Dzaakk/simple-commerce/internal/shopping_cart/models"
	template "Dzaakk/simple-commerce/package/templates"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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
	queryCreateProduct    = `INSERT INTO public.product (product_name, price, stock, category_id, seller_id, created, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	queryUpdate           = `UPDATE public.product SET product_name=$1, price=$2, stock=$3, updated=NOW(), updated_by=$4 WHERE id=$5`
	queryFindByCategoryId = `SELECT * FROM public.product WHERE category_id = $1`
	queryFindBySellerId   = `SELECT * FROM public.product WHERE seller_id = $1`
	queryFindById         = `SELECT * FROM public.product WHERE id = $1`
	queryGetPriceById     = `SELECT price FROM public.product WHERE id = $1`
	queryGetStockById     = `SELECT stock FROM public.product WHERE id = $1`
	queryFindByName       = `SELECT * FROM public.product WHERE product_name like '%' || $1 || '%'`
	querySetStockById     = `UPDATE public.product SET stock = $1 WHERE id = $2`
)

func (repo *ProductRepositoryImpl) Create(data model.TProduct) (*model.TProduct, error) {
	statement, err := repo.DB.Prepare(queryCreateProduct)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var id int

	err = statement.QueryRow(data.ProductName, data.Price, data.Stock, data.CategoryId, data.SellerId, data.Base.Created, data.Base.CreatedBy).Scan(id)
	if err != nil {
		return nil, err
	}

	data.Id = id
	return &data, nil
}

func (repo *ProductRepositoryImpl) Update(data model.TProduct) (int64, error) {
	result, err := repo.DB.Exec(
		queryUpdate,
		data.ProductName,
		data.Price,
		data.Stock,
		data.UpdatedBy,
		data.Id)

	if err != nil {
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

func (repo *ProductRepositoryImpl) FindBySellerId(sellerId int) ([]*model.TProduct, error) {
	rows, err := repo.DB.Query(queryFindBySellerId, sellerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listProduct []*model.TProduct
	for rows.Next() {
		product, err := retrieveProduct(rows)
		if err != nil {
			return nil, err
		}
		listProduct = append(listProduct, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listProduct, nil
}

func (repo *ProductRepositoryImpl) FindBySellerIdAndCategoryId(sellerId int, categoryId int) ([]*model.TProduct, error) {
	panic("unimplemented")
}

func (repo *ProductRepositoryImpl) FindByCategoryId(categoryId int) ([]*model.TProduct, error) {
	rows, err := repo.DB.Query(queryFindByCategoryId, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listProduct []*model.TProduct
	for rows.Next() {
		product, err := retrieveProduct(rows)
		if err != nil {
			return nil, err
		}
		listProduct = append(listProduct, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listProduct, nil
}

func (repo *ProductRepositoryImpl) FindById(id int) (*model.TProduct, error) {
	rows, err := repo.DB.Query(queryFindById, id)
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

func (repo *ProductRepositoryImpl) GetPriceById(id int) (*float32, error) {
	var balance float32
	err := repo.DB.QueryRow(queryGetPriceById, id).Scan(&balance)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (repo *ProductRepositoryImpl) GetStockById(id int) (int, error) {
	var stock int
	err := repo.DB.QueryRow(queryGetPriceById, id).Scan(stock)
	if err != nil {
		return 0, err
	}

	return stock, nil
}

func (repo *ProductRepositoryImpl) SetStockById(id int, stock int) (int64, error) {
	result, err := repo.DB.Exec(querySetStockById, stock, id)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

func (repo *ProductRepositoryImpl) FindByName(name string) (*model.TProduct, error) {

	rows, err := repo.DB.Query(queryFindByName, name)
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

func (repo *ProductRepositoryImpl) UpdateStock(listData []*models.TCartItemDetail, name string) error {
	query, args := generateMultipleStockUpdateQuery(listData)
	_, err := repo.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductRepositoryImpl) UpdateStockWithTx(tx *sql.Tx, listItem []*models.TCartItemDetail) ([]*int, error) {
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

func verifyStockAvailability(tx *sql.Tx, listItem []*models.TCartItemDetail) ([]*int, error) {
	var query strings.Builder
	var args []interface{}
	listEmptyProductId := []*int{}

	query.WriteString(`SELECT id, stock FROM public.product WHERE id IN (`)

	for i, item := range listItem {
		if i > 0 {
			query.WriteString(",")
		}
		query.WriteString(fmt.Sprintf("$%d", i+1))
		args = append(args, item.ProductId)
	}

	query.WriteString(") FOR UPDATE")

	rows, err := tx.Query(query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("error locking products: %w", err)
	}
	defer rows.Close()

	stockMap := make(map[int]int)
	for rows.Next() {
		var id, stock int
		if err := rows.Scan(&id, &stock); err != nil {
			return nil, fmt.Errorf("error scanning product stock: %w", err)
		}
		stockMap[id] = stock
	}

	for _, item := range listItem {
		currentStock, exists := stockMap[item.ProductId]
		if !exists {
			return nil, fmt.Errorf("product %d not found", item.ProductId)
		}
		if currentStock == 1 {
			listEmptyProductId = append(listEmptyProductId, &item.ProductId)
		}
		if currentStock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %d: requested %d, available %d", item.ProductId, item.Quantity, currentStock)
		}
	}
	return listEmptyProductId, nil
}

func generateMultipleStockUpdateQuery(listData []*models.TCartItemDetail) (string, []interface{}) {
	var query strings.Builder
	var args []interface{}
	query.WriteString("UPDATE public.product SET stock = CASE id ")

	for _, item := range listData {
		query.WriteString(fmt.Sprintf("WHEN $%d THEN stock - $%d ", len(args)+1, len(args)+2))
		args = append(args, item.ProductId, item.Quantity)
	}

	query.WriteString(" END, updated_by = 'SYSTEM', updated = NOW() WHERE id IN (")

	for i, item := range listData {
		if i > 0 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("$%d", len(args)+1))
		args = append(args, item.ProductId)
	}

	query.WriteString(")")
	return query.String(), args
}

func rowsToProduct(rows *sql.Rows) (*model.TProduct, error) {
	base := template.Base{}
	product := model.TProduct{}

	err := rows.Scan(&product.Id, &product.ProductName, &product.Price, &product.Stock, &product.CategoryId, &base.Created, &base.CreatedBy, &base.Updated, &base.UpdatedBy)
	if err != nil {
		return nil, err
	}
	if !base.UpdatedBy.Valid {
		base.UpdatedBy.String = ""
	}
	product.Base = base

	return &product, nil
}

func retrieveProduct(rows *sql.Rows) (*model.TProduct, error) {
	if rows.Next() {
		return rowsToProduct(rows)
	}
	return nil, errors.New("product not found")
}
