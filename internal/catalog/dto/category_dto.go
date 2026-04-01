package dto

import (
	"Dzaakk/simple-commerce/internal/catalog/model"
	"time"
)

type CategoryTree struct {
	ID       int64
	ParentID *int64
	Name     string
	Slug     string
	Depth    int
}

type CreateCategoryReq struct {
	ParentID *int64 `json:"parent_id"`
	Name     string `json:"name" validate:"required"`
	Slug     string `json:"slug" validate:"required"`
	IsActive *bool  `json:"is_active"`
}

func (c *CreateCategoryReq) ToCreateData() *model.Category {
	isActive := true
	if c.IsActive != nil {
		isActive = *c.IsActive
	}

	return &model.Category{
		ParentID:  c.ParentID,
		Name:      c.Name,
		Slug:      c.Slug,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
