package service

import (
	"Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/internal/catalog/model"
	repo "Dzaakk/simple-commerce/internal/catalog/repository"
	"context"
	"time"
)

type ProductService interface {
	Create(ctx context.Context, req *dto.CreateProductReq) (string, error)
	Update(ctx context.Context, productID string, sellerID string, req *dto.UpdateProductReq) error
	SoftDelete(ctx context.Context, productID string, sellerID string) error
	FindByID(ctx context.Context, productID string) (*dto.ProductRes, error)
	FindAll(ctx context.Context, req dto.ProductQueryReq) (*dto.ProductListRes, error)
}

type ProductRepository interface {
	Create(ctx context.Context, data *model.Product) (string, error)
	Update(ctx context.Context, data *model.Product) (int64, error)
	SoftDelete(ctx context.Context, productID string, updatedAt time.Time) (int64, error)
	FindByID(ctx context.Context, productID string) (*model.Product, error)
	FindBySellerID(ctx context.Context, sellerID string) ([]*model.Product, error)
	FindAll(ctx context.Context, filter repo.ProductFilter) ([]*model.Product, error)
}

type CategoryService interface {
	Create(ctx context.Context, req *dto.CreateCategoryReq) (int64, error)
	FindAll(ctx context.Context) ([]*dto.CategoryTree, error)
	FindByID(ctx context.Context, categoryID int64) (*dto.CategoryTree, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, data *model.Category) (int64, error)
	FindByID(ctx context.Context, id int64) (*model.Category, error)
	FindAll(ctx context.Context) ([]*dto.CategoryTree, error)
}
