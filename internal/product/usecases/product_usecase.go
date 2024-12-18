package usecases

import (
	model "Dzaakk/simple-commerce/internal/product/models"
)

type ProductUseCase interface {
	FindByCategoryId(categoryId int) ([]*model.ProductRes, error)
	Create(data model.ProductReq) (*model.ProductRes, error)
	Update(data model.ProductReq) error
	FilterByPrice(price int) ([]*model.ProductRes, error)
	FindByName(productName string) (*model.ProductRes, error)
}
