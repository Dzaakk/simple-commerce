package dto

import (
	"Dzaakk/simple-commerce/internal/catalog/model"
	"time"
)

type CreateProductReq struct {
	SellerID    string  `json:"seller_id" validate:"required"`
	CategoryID  int64   `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	SKU         string  `json:"sku" validate:"required"`
	Description *string `json:"description"`
	Price       float64 `json:"price" validate:"required"`
	ImageURL    *string `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
}

func (c *CreateProductReq) ToCreateData() *model.Product {
	isActive := true
	if c.IsActive != nil {
		isActive = *c.IsActive
	}

	return &model.Product{
		SellerID:    c.SellerID,
		CategoryID:  c.CategoryID,
		Name:        c.Name,
		SKU:         c.SKU,
		Description: c.Description,
		Price:       c.Price,
		ImageURL:    c.ImageURL,
		IsActive:    isActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

type UpdateProductReq struct {
	CategoryID  int64   `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	SKU         string  `json:"sku" validate:"required"`
	Description *string `json:"description"`
	Price       float64 `json:"price" validate:"required"`
	ImageURL    *string `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
}

func (u *UpdateProductReq) ToUpdateData(productID string, sellerID string) *model.Product {
	isActive := true
	if u.IsActive != nil {
		isActive = *u.IsActive
	}

	return &model.Product{
		ID:          productID,
		SellerID:    sellerID,
		CategoryID:  u.CategoryID,
		Name:        u.Name,
		SKU:         u.SKU,
		Description: u.Description,
		Price:       u.Price,
		ImageURL:    u.ImageURL,
		IsActive:    isActive,
		UpdatedAt:   time.Now(),
	}
}

type ProductRes struct {
	ID          string    `json:"id"`
	SellerID    string    `json:"seller_id,omitempty"`
	CategoryID  int64     `json:"category_id,omitempty"`
	Name        string    `json:"name,omitempty"`
	SKU         string    `json:"sku,omitempty"`
	Description *string   `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	ImageURL    *string   `json:"image_url,omitempty"`
	IsActive    bool      `json:"is_active,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func ToProductRes(p *model.Product) ProductRes {
	return ProductRes{
		ID:          p.ID,
		SellerID:    p.SellerID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		SKU:         p.SKU,
		Description: p.Description,
		Price:       p.Price,
		ImageURL:    p.ImageURL,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type ProductListRes struct {
	Items      []ProductRes `json:"items"`
	NextCursor *string      `json:"next_cursor,omitempty"`
}

type ProductQueryReq struct {
	CategoryID *int64   `json:"category_id"`
	SellerID   *string  `json:"seller_id"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	Name       *string  `json:"name"`
	Cursor     *string  `json:"cursor"`
	Limit      int      `json:"limit"`
	SortBy     string   `json:"sort_by"` // "price_asc", "price_desc", "newest"
}
