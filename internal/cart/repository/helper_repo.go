package repository

import (
	"Dzaakk/simple-commerce/internal/cart/model"
	"Dzaakk/simple-commerce/package/response"
	"database/sql"
	"errors"
)

func scanCart(row *sql.Row) (*model.Cart, error) {
	var cart model.Cart

	err := row.Scan(
		&cart.ID,
		&cart.CustomerID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, response.Error("failed to scan cart", err)
	}

	return &cart, nil
}
