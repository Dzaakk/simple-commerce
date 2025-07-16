package usecase

import (
	"Dzaakk/simple-commerce/internal/product/model"
	"context"
)

type ProductUseCase interface {
	FindByCategoryID(ctx context.Context, categoryID int) ([]*model.ProductRes, error)
	Create(ctx context.Context, data model.ProductReq) (*model.ProductRes, error)
	Update(ctx context.Context, data model.ProductReq) error
	FindByProductName(ctx context.Context, productName string) (*model.ProductRes, error)
	FindByFilter(ctx context.Context, params model.ProductFilter) ([]*model.ProductRes, error)
}
