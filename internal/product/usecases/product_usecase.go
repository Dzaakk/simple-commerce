package usecases

import (
	model "Dzaakk/simple-commerce/internal/product/models"
)

type ProductUseCase interface {
	FindByCategoryId(categoryId int) ([]*model.ProductRes, error)
	Create(data model.TProduct) (*model.ProductRes, error)
	Update(data model.TProduct) error
	FilterByPrice(price int) ([]*model.ProductRes, error)
	FindByName(productName string) (*model.ProductRes, error)
}
