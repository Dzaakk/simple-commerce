package repository

import (
	"Dzaakk/simple-commerce/internal/catalog/model"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
)

const (
	categorySelectColumns = "id, parent_id, name, slug, is_active, created_at, updated_at"
	categoryQueryCreate   = "INSERT INTO public.categories (parent_id, name, slug, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	categoryQueryFindByID = "SELECT " + categorySelectColumns + " FROM public.categories WHERE id=$1"
	categoryQueryFindAll  = "SELECT " + categorySelectColumns + " FROM public.categories WHERE is_active = true ORDER BY id ASC"
)

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) Create(ctx context.Context, data *model.Category) (int64, error) {
	var id int64

	err := r.DB.QueryRowContext(
		ctx,
		categoryQueryCreate,
		data.ParentID,
		data.Name,
		data.Slug,
		data.IsActive,
		data.CreatedAt,
		data.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, response.Error("failed to create category", err)
	}

	return id, nil
}

func (r *CategoryRepository) FindByID(ctx context.Context, id int64) (*model.Category, error) {
	row := r.DB.QueryRowContext(ctx, categoryQueryFindByID, id)

	return scanCategory(row)
}

func (r *CategoryRepository) FindAll(ctx context.Context) ([]*model.Category, error) {
	rows, err := r.DB.QueryContext(ctx, categoryQueryFindAll)
	if err != nil {
		return nil, response.Error("failed to query categories", err)
	}
	defer rows.Close()

	var categories []*model.Category

	for rows.Next() {
		var c model.Category
		err := rows.Scan(
			&c.ID,
			&c.ParentID,
			&c.Name,
			&c.Slug,
			&c.IsActive,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, response.Error("failed to scan category", err)
		}

		categories = append(categories, &c)
	}

	return categories, nil
}
