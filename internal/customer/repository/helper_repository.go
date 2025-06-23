package repository

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	"database/sql"
	"errors"
	"fmt"
)

func scanCustomer(row *sql.Row) (*model.TCustomers, error) {
	customer := &model.TCustomers{}
	var updated sql.NullTime

	err := row.Scan(
		&customer.ID, &customer.Username, &customer.Email, &customer.Password, &customer.PhoneNumber, &customer.Balance, &customer.Status,
		&customer.Created, &customer.CreatedBy, &updated, &customer.UpdatedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan customer: %w", err)
	}

	if updated.Valid {
		customer.Updated.Time = updated.Time
	}
	if !customer.UpdatedBy.Valid {
		customer.UpdatedBy.String = ""
	}

	return customer, nil
}
