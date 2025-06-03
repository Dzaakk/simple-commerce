package usecase

import (
	"Dzaakk/simple-commerce/internal/product/model"
	"context"
)

type ProductUseCase interface {
	FindByCategoryID(ctx context.Context, categoryID int) ([]*model.ProductRes, error)
	Create(ctx context.Context, data model.ProductReq) (*model.ProductRes, error)
	Update(ctx context.Context, data model.ProductReq) error
	FilterByProductPrice(ctx context.Context, productPrice int) ([]*model.ProductRes, error)
	FindByProductName(ctx context.Context, productName string) (*model.ProductRes, error)
}
