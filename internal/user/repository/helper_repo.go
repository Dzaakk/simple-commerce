package repository

import (
	"Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/response"
	"database/sql"
	"errors"
)

func scanCustomer(row *sql.Row) (*model.Customer, error) {
	customer := &model.Customer{}

	err := row.Scan(
		&customer.ID,
		&customer.Email,
		&customer.PasswordHash,
		&customer.FullName,
		&customer.Phone,
		&customer.Status,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, response.Error("failed to scan customer", err)
	}

	return customer, nil
}

func scanSeller(row *sql.Row) (*model.Seller, error) {
	seller := &model.Seller{}

	err := row.Scan(
		&seller.ID,
		&seller.Email,
		&seller.PasswordHash,
		&seller.ShopName,
		&seller.Phone,
		&seller.Status,
		&seller.CreatedAt,
		&seller.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, response.Error("failed to scan seller", err)
	}

	return seller, nil
}
