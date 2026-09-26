package service

import (
	"context"

	"Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/internal/catalog/model"
)

type ProductService interface {
	FindByID(context.Context, int64) (*dto.ProductResponse, error)
	FindByIDCached(context.Context, int64) (*dto.ProductResponse, error)
	FindAll(context.Context, dto.ProductQuery) (*dto.ProductListResponse, error)
	FindAllCached(context.Context, dto.ProductQuery) (*dto.ProductListResponse, error)
}

type ProductRepository interface {
	FindByID(context.Context, int64) (*model.Product, error)
	FindAll(context.Context, dto.ProductQuery) ([]*model.Product, error)
}

type CategoryService interface {
	FindAll(context.Context) ([]dto.CategoryResponse, error)
}

type CategoryRepository interface {
	FindAll(context.Context) ([]*model.Category, error)
}
