package repository

import (
	"context"
	"database/sql"

	"Dzaakk/simple-commerce/internal/catalog/model"
	"Dzaakk/simple-commerce/package/response"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll(ctx context.Context) ([]*model.Category, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, slug, created_at FROM categories ORDER BY id")
	if err != nil {
		return nil, response.Error("failed to list categories", err)
	}
	defer rows.Close()

	categories := make([]*model.Category, 0)
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Slug, &category.CreatedAt); err != nil {
			return nil, response.Error("failed to scan category", err)
		}
		categories = append(categories, &category)
	}
	if err := rows.Err(); err != nil {
		return nil, response.Error("failed to iterate categories", err)
	}
	return categories, nil
}
