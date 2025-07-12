package repository

import (
	"Dzaakk/simple-commerce/internal/seller/model"
	"Dzaakk/simple-commerce/package/template"
	"database/sql"
	"errors"
	"fmt"
)

func scanSeller(row *sql.Row) (*model.TSeller, error) {
	seller := &model.TSeller{}
	base := template.Base{}
	var updated sql.NullTime

	err := row.Scan(
		&seller.ID, &seller.Username, &seller.Email, &seller.Password, &seller.Balance, &seller.Status,
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

func scanListSeller(rows *sql.Rows) ([]*model.TSeller, error) {
	var listSeller []*model.TSeller

	for rows.Next() {
		seller := &model.TSeller{}
		base := template.Base{}
		var updated sql.NullTime

		err := rows.Scan(
			&seller.ID, &seller.Username, &seller.Email, &seller.Password, &seller.Balance, &seller.Status,
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
		listSeller = append(listSeller, seller)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return listSeller, nil
}
