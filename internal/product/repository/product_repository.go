package repository

import (
	model "Dzaakk/simple-commerce/internal/product/models"
)

type ProductRepository interface {
	Create(data model.TProduct) (*model.ProductRes, error)
	Update(data model.TProduct) (*model.ProductRes, error)
	FindByCategoryId(categoryId int) ([]*model.TProduct, error)
	FindById(id int) (*model.TProduct, error)
	GetPriceById(id int) (*float32, error)
	GetStockById(id int) (int, error)
}
