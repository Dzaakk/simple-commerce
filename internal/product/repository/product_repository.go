package product

import (
	model "Dzaakk/synapsis/internal/product/models"
	"Dzaakk/synapsis/package/template"
	"database/sql"
	"errors"
	"fmt"
)

type ProductRepository interface {
	Create(data model.TProduct) (*model.ProductRes, error)
	Update(data model.TProduct) (*model.ProductRes, error)
	FindByCategoryId(categoryId int) ([]*model.TProduct, error)
	FindById(id int) (*model.TProduct, error)
	GetPriceById(id int) (*float32, error)
	GetStockById(id int) (int, error)
}

type ProductRepositoryImpl struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &ProductRepositoryImpl{
		DB: db,
	}
}

const queryCreateProduct = `INSERT INTO public.product (product_name, price, stock, category_id, created, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

func (repo *ProductRepositoryImpl) Create(data model.TProduct) (*model.ProductRes, error) {
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

	newProduct := &model.ProductRes{
		Id:          fmt.Sprintf("%d", id),
		ProductName: data.ProductName,
		Price:       fmt.Sprintf("%0.f", data.Price),
		Stock:       fmt.Sprintf("%d", data.Stock),
		CategoryId:  fmt.Sprintf("%d", data.CategoryId),
	}
	return newProduct, nil
}

func (repo *ProductRepositoryImpl) Update(data model.TProduct) (*model.ProductRes, error) {
	panic("unimplemented")
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
