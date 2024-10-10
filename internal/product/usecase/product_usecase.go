package usecase

import (
	model "Dzaakk/simple-commerce/internal/product/models"
)

type ProductUseCase interface {
	FindByCategoryId(categoryId int) ([]*model.ProductRes, error)
}
