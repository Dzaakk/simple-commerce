package repository

import (
	"Dzaakk/simple-commerce/internal/user/domain"
	"Dzaakk/simple-commerce/package/response"
	"database/sql"
	"errors"
)

func scanCustomer(row *sql.Row) (*domain.Customer, error) {
	customer := &domain.Customer{}

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
