package dto

import (
	"time"

	"Dzaakk/simple-commerce/internal/catalog/model"
)

type ProductResponse struct {
	ID          int64     `json:"id"`
	CategoryID  int64     `json:"category_id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToProductResponse(product *model.Product) ProductResponse {
	return ProductResponse{
		ID: product.ID, CategoryID: product.CategoryID, SKU: product.SKU,
		Name: product.Name, Description: product.Description, Price: product.Price,
		IsActive: product.IsActive, CreatedAt: product.CreatedAt, UpdatedAt: product.UpdatedAt,
	}
}

type ProductListResponse struct {
	Items      []ProductResponse `json:"items"`
	NextCursor *string           `json:"next_cursor,omitempty"`
}

type ProductQuery struct {
	CategoryID *int64
	MinPrice   *float64
	MaxPrice   *float64
	Name       *string
	Cursor     *string
	Limit      int
	SortBy     string
}
