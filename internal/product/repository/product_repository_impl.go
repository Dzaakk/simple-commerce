package repository

import (
	model "Dzaakk/simple-commerce/internal/product/models"
	"Dzaakk/simple-commerce/package/template"
	"database/sql"
	"errors"
	"fmt"
)

type ProductRepositoryImpl struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &ProductRepositoryImpl{
		DB: db,
	}
}

const queryCreateProduct = `INSERT INTO public.product (product_name, price, stock, category_id, created, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

func (repo *ProductRepositoryImpl) Create(data model.TProduct) (*model.TProduct, error) {
	statement, err := repo.DB.Prepare(queryCreateProduct)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var id int

	err = statement.QueryRow(data.ProductName, data.Price, data.Stock, data.CategoryId, data.Base.Created, data.Base.CreatedBy).Scan(id)
	if err != nil {
		return nil, err
	}

	data.Id = id
	return &data, nil
}

const queryUpdate = `UPDATE public.product SET product_name=$1, price=$2, stock=$3, updated=NOW(), updated_by=$4 WHERE id=$5`

func (repo *ProductRepositoryImpl) Update(data model.TProduct) error {
	statement, err := repo.DB.Prepare(queryUpdate)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(data.ProductName, data.Price, data.Stock, data.UpdatedBy, data.Id)
	if err != nil {
		return err
	}

	return nil
}

const queryFindByCategoryId = `SELECT * FROM public.product WHERE category_id = $1`

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
	fmt.Println("LEN = ", len(listProduct))
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listProduct, nil
}

const queryFindById = `SELECT * FROM public.product WHERE id = $1`

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

const queryGetPriceById = `SELECT price FROM public.product WHERE id = $1`

func (repo *ProductRepositoryImpl) GetPriceById(id int) (*float32, error) {
	var balance float32
	err := repo.DB.QueryRow(queryGetPriceById, id).Scan(&balance)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

const queryGetStockById = `SELECT stock FROM public.product WHERE id = $1`

func (repo *ProductRepositoryImpl) GetStockById(id int) (int, error) {
	var stock int
	err := repo.DB.QueryRow(queryGetPriceById, id).Scan(stock)
	if err != nil {
		return 0, err
	}

	return stock, nil
}

const queryFindByName = `SELECT * FROM public.product WHERE product_name like '%' || $1 || '%'`

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
