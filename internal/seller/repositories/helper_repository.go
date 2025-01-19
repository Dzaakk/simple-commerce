package repositories

import (
	model "Dzaakk/simple-commerce/internal/seller/models"
	template "Dzaakk/simple-commerce/package/templates"
	"database/sql"
	"errors"
	"fmt"
)

func scanSeller(row *sql.Row) (*model.TSeller, error) {
	seller := &model.TSeller{}
	base := template.Base{}
	var updated sql.NullTime

	err := row.Scan(
		&seller.Id, &seller.Username, &seller.Email, &seller.Password, &seller.Balance, &seller.Status,
		&base.Created, &base.CreatedBy, &updated, &base.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan seller: %w", err)
	}

	if updated.Valid {
		base.Updated.Time = updated.Time
	}
	if !base.UpdatedBy.Valid {
		base.UpdatedBy.String = ""
	}

	seller.Base = base

	return seller, nil
}