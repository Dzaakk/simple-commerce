package repository

import (
	"Dzaakk/simple-commerce/internal/catalog/model"
	"Dzaakk/simple-commerce/package/response"
	"database/sql"
)

func scanCategory(row *sql.Row) (*model.Category, error) {
	var c model.Category

	err := row.Scan(
		&c.ID,
		&c.ParentID,
		&c.Name,
		&c.Slug,
		&c.IsActive,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, response.Error("category not found", err)
		}
		return nil, response.Error("failed to scan category", err)
	}

	return &c, nil
}

func scanProduct(row *sql.Row) (*model.Product, error) {
	var p model.Product

	err := row.Scan(
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
		if err == sql.ErrNoRows {
			return nil, response.Error("product not found", err)
		}
		return nil, response.Error("failed to scan product", err)
	}

	return &p, nil
}
