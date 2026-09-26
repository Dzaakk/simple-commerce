package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/internal/catalog/model"
	"Dzaakk/simple-commerce/package/response"
)

const productColumns = "id, category_id, sku, name, description, price, is_active, created_at, updated_at"

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindByID(ctx context.Context, id int64) (*model.Product, error) {
	row := r.db.QueryRowContext(ctx, "SELECT "+productColumns+" FROM products WHERE id = $1 AND is_active = true", id)
	product, err := scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, response.Error("failed to find product", err)
	}
	return product, nil
}

func (r *ProductRepository) FindAll(ctx context.Context, filter dto.ProductQuery) ([]*model.Product, error) {
	query, args := buildProductQuery(filter)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, response.Error("failed to list products", err)
	}
	defer rows.Close()

	products := make([]*model.Product, 0)
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, response.Error("failed to scan product", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, response.Error("failed to iterate products", err)
	}
	return products, nil
}

type rowScanner interface {
	Scan(...any) error
}

func scanProduct(row rowScanner) (*model.Product, error) {
	var product model.Product
	err := row.Scan(
		&product.ID, &product.CategoryID, &product.SKU, &product.Name,
		&product.Description, &product.Price, &product.IsActive,
		&product.CreatedAt, &product.UpdatedAt,
	)
	return &product, err
}

func buildProductQuery(filter dto.ProductQuery) (string, []any) {
	query := "SELECT " + productColumns + " FROM products WHERE is_active = true"
	args := make([]any, 0, 8)
	add := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}

	if filter.CategoryID != nil {
		query += " AND category_id = " + add(*filter.CategoryID)
	}
	if filter.MinPrice != nil {
		query += " AND price >= " + add(*filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query += " AND price <= " + add(*filter.MaxPrice)
	}
	if filter.Name != nil {
		query += " AND name ILIKE " + add("%"+*filter.Name+"%")
	}

	if filter.Cursor != nil {
		value, id, ok := splitCursor(*filter.Cursor)
		if ok {
			switch filter.SortBy {
			case "price_asc":
				if price, err := strconv.ParseFloat(value, 64); err == nil {
					query += " AND (price, id) > (" + add(price) + ", " + add(id) + ")"
				}
			case "price_desc":
				if price, err := strconv.ParseFloat(value, 64); err == nil {
					query += " AND (price, id) < (" + add(price) + ", " + add(id) + ")"
				}
			default:
				if createdAt, err := time.Parse(time.RFC3339Nano, value); err == nil {
					query += " AND (created_at, id) < (" + add(createdAt) + ", " + add(id) + ")"
				}
			}
		}
	}

	switch filter.SortBy {
	case "price_asc":
		query += " ORDER BY price ASC, id ASC"
	case "price_desc":
		query += " ORDER BY price DESC, id DESC"
	default:
		query += " ORDER BY created_at DESC, id DESC"
	}
	query += " LIMIT " + add(filter.Limit)
	return query, args
}

func splitCursor(cursor string) (string, int64, bool) {
	parts := strings.SplitN(cursor, "|", 2)
	if len(parts) != 2 {
		return "", 0, false
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	return parts[0], id, err == nil
}
