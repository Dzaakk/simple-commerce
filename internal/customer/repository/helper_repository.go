package repositories

import (
	model "Dzaakk/simple-commerce/internal/customer/models"
	template "Dzaakk/simple-commerce/package/templates"
	"database/sql"
	"errors"
	"fmt"
)

func scanCustomer(row *sql.Row) (*model.TCustomers, error) {
	customer := &model.TCustomers{}
	base := template.Base{}
	var updated sql.NullTime

	err := row.Scan(
		&customer.ID, &customer.Username, &customer.Email, &customer.Password, &customer.PhoneNumber, &customer.Balance, &customer.Status,
		&base.Created, &base.CreatedBy, &updated, &base.UpdatedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan customer: %w", err)
	}
	if updated.Valid {
		base.Updated.Time = updated.Time
	}
	if !base.UpdatedBy.Valid {
		base.UpdatedBy.String = ""
	}

	customer.Base = base

	return customer, nil
}
