package usecase

import (
	model "Dzaakk/simple-commerce/internal/product/models"
)

type ProductUseCase interface {
	FindByCategoryId(categoryId int) ([]*model.ProductRes, error)
	Create(data model.TProduct) (*int, error)
	UpdateStock(id, stock int) (*int, error)
	Update(data model.TProduct) error
	FilterByPrice(price int) ([]*model.ProductRes, error)
}
