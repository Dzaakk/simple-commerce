package repository

import (
	"Dzaakk/simple-commerce/internal/catalog/model"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	productSelectColumns     = "id, seller_id, category_id, name, sku, description, price, image_url, is_active, created_at, updated_at"
	productQueryCreate       = "INSERT INTO public.products (seller_id, category_id, name, sku, description, price, image_url, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id"
	productQueryUpdate       = "UPDATE public.products SET seller_id=$1, category_id=$2, name=$3, sku=$4, description=$5, price=$6, image_url=$7, is_active=$8, updated_at=$9 WHERE id=$10"
	productQuerySoftDelete   = "UPDATE public.products SET is_active=false, updated_at=$1 WHERE id=$2"
	productQueryFindByID     = "SELECT " + productSelectColumns + " FROM public.products WHERE id=$1"
	productQueryFindBySeller = "SELECT " + productSelectColumns + " FROM public.products WHERE seller_id=$1 AND is_active=true"
)

type ProductFilter struct {
	CategoryID *int64
	SellerID   *string
	MinPrice   *float64
	MaxPrice   *float64
	Name       *string // search by name (ILIKE)
	Cursor     *string // pagination cursor: "value|id"
	Limit      int
	SortBy     string // "price_asc", "price_desc", "newest"
}

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) Create(ctx context.Context, data *model.Product) (string, error) {
	var id string

	err := r.DB.QueryRowContext(
		ctx,
		productQueryCreate,
		data.SellerID,
		data.CategoryID,
		data.Name,
		data.SKU,
		data.Description,
		data.Price,
		data.ImageURL,
		data.IsActive,
		data.CreatedAt,
		data.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return "", response.Error("failed to create product", err)
	}

	return id, nil
}

func (r *ProductRepository) Update(ctx context.Context, data *model.Product) (int64, error) {
	result, err := r.DB.ExecContext(
		ctx,
		productQueryUpdate,
		data.SellerID,
		data.CategoryID,
		data.Name,
		data.SKU,
		data.Description,
		data.Price,
		data.ImageURL,
		data.IsActive,
		data.UpdatedAt,
		data.ID,
	)

	if err != nil {
		return 0, response.ExecError("update product", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, response.Error("failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return 0, response.Error("no rows updated", sql.ErrNoRows)
	}

	return rowsAffected, nil
}

func (r *ProductRepository) SoftDelete(ctx context.Context, id string, updatedAt time.Time) (int64, error) {
	result, err := r.DB.ExecContext(ctx, productQuerySoftDelete, updatedAt, id)
	if err != nil {
		return 0, response.ExecError("soft delete product", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, response.Error("failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return 0, response.Error("no rows updated", sql.ErrNoRows)
	}

	return rowsAffected, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*model.Product, error) {
	row := r.DB.QueryRowContext(ctx, productQueryFindByID, id)

	return scanProduct(row)
}

func (r *ProductRepository) FindBySellerID(ctx context.Context, sellerID string) ([]*model.Product, error) {
	rows, err := r.DB.QueryContext(ctx, productQueryFindBySeller, sellerID)
	if err != nil {
		return nil, response.Error("failed to query products by seller", err)
	}
	defer rows.Close()

	var products []*model.Product

	for rows.Next() {
		var p model.Product
		err := rows.Scan(
			&p.ID,
			&p.SellerID,
			&p.CategoryID,
			&p.Name,
			&p.SKU,
			&p.Description,
			&p.Price,
			&p.ImageURL,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, response.Error("failed to scan product", err)
		}

		products = append(products, &p)
	}

	return products, nil
}

func (r *ProductRepository) FindAll(ctx context.Context, filter ProductFilter) ([]*model.Product, error) {
	query, args := buildProductQuery(filter)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, response.Error("failed to query products", err)
	}
	defer rows.Close()

	var products []*model.Product

	for rows.Next() {
		var p model.Product
		err := rows.Scan(
			&p.ID,
			&p.SellerID,
			&p.CategoryID,
			&p.Name,
			&p.SKU,
			&p.Description,
			&p.Price,
			&p.ImageURL,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, response.Error("failed to scan product", err)
		}

		products = append(products, &p)
	}

	return products, nil
}

func buildProductQuery(f ProductFilter) (string, []any) {
	query := "SELECT " + productSelectColumns + " FROM public.products WHERE is_active = true"
	args := []any{}
	argPos := 1

	if f.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argPos)
		args = append(args, *f.CategoryID)
		argPos++
	}
	if f.SellerID != nil {
		query += fmt.Sprintf(" AND seller_id = $%d", argPos)
		args = append(args, *f.SellerID)
		argPos++
	}
	if f.MinPrice != nil {
		query += fmt.Sprintf(" AND price >= $%d", argPos)
		args = append(args, *f.MinPrice)
		argPos++
	}
	if f.MaxPrice != nil {
		query += fmt.Sprintf(" AND price <= $%d", argPos)
		args = append(args, *f.MaxPrice)
		argPos++
	}
	if f.Name != nil {
		query += fmt.Sprintf(" AND name ILIKE $%d", argPos)
		args = append(args, "%"+*f.Name+"%")
		argPos++
	}

	sortBy := f.SortBy
	if sortBy == "" {
		sortBy = "newest"
	}

	if f.Cursor != nil && *f.Cursor != "" {
		cursorVal, cursorID, hasID := splitCursor(*f.Cursor)

		switch sortBy {
		case "price_asc":
			if hasID {
				if price, err := strconv.ParseFloat(cursorVal, 64); err == nil {
					query += fmt.Sprintf(" AND (price, id) > ($%d, $%d)", argPos, argPos+1)
					args = append(args, price, cursorID)
					argPos += 2
				}
			}
		case "price_desc":
			if hasID {
				if price, err := strconv.ParseFloat(cursorVal, 64); err == nil {
					query += fmt.Sprintf(" AND (price, id) < ($%d, $%d)", argPos, argPos+1)
					args = append(args, price, cursorID)
					argPos += 2
				}
			}
		default: // newest
			if t, err := time.Parse(time.RFC3339Nano, cursorVal); err == nil {
				if hasID {
					query += fmt.Sprintf(" AND (created_at, id) < ($%d, $%d)", argPos, argPos+1)
					args = append(args, t, cursorID)
					argPos += 2
				} else {
					query += fmt.Sprintf(" AND created_at < $%d", argPos)
					args = append(args, t)
					argPos++
				}
			}
		}
	}

	switch sortBy {
	case "price_asc":
		query += " ORDER BY price ASC, id ASC"
	case "price_desc":
		query += " ORDER BY price DESC, id DESC"
	default: // newest
		query += " ORDER BY created_at DESC, id DESC"
	}

	if f.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, f.Limit)
		argPos++
	}

	return query, args
}

func splitCursor(cursor string) (string, string, bool) {
	parts := strings.SplitN(cursor, "|", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], true
	}
	return cursor, "", false
}
