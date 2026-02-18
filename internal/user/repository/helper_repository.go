package repository

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	"Dzaakk/simple-commerce/package/response"
	"database/sql"
	"errors"
)

func scanCustomer(row *sql.Row) (*model.Customers, error) {
	customer := &model.Customers{}
	var updated sql.NullTime

	err := row.Scan(
		&customer.ID,
		&customer.Username,
		&customer.Email,
		&customer.Password,
		&customer.Gender,
		&customer.PhoneNumber,
		&customer.Balance,
		&customer.Status,
		&customer.DateOfBirth,
		&customer.ProfilePicture,
		&customer.LastLogin,
		&customer.Created,
		&customer.CreatedBy,
		&updated,
		&customer.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, response.Error("failed to scan customer", err)
	}

	if updated.Valid {
		customer.Updated = updated
	}

	return customer, nil
}
