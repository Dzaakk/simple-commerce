package dto

import (
	"time"

	"Dzaakk/simple-commerce/internal/catalog/model"
)

type CategoryResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

func ToCategoryResponse(category *model.Category) CategoryResponse {
	return CategoryResponse{ID: category.ID, Name: category.Name, Slug: category.Slug, CreatedAt: category.CreatedAt}
}
