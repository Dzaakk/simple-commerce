package usecases

import (
	model "Dzaakk/simple-commerce/internal/product/models"
	"context"
)

type ProductUseCase interface {
	FindByCategoryId(ctx context.Context, categoryId int) ([]*model.ProductRes, error)
	Create(ctx context.Context, data model.ProductReq) (*model.ProductRes, error)
	Update(ctx context.Context, data model.ProductReq) error
	FilterByPrice(ctx context.Context, price int) ([]*model.ProductRes, error)
	FindByName(ctx context.Context, productName string) (*model.ProductRes, error)
}
