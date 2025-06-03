package repository

import (
	"Dzaakk/simple-commerce/internal/shopping_cart/model"
	response "Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/template"
	"database/sql"
	"errors"
)

func scanCart(row *sql.Row) (*model.TShoppingCart, error) {
	cart := &model.TShoppingCart{}
	base := template.Base{}
	var updated sql.NullTime

	err := row.Scan(
		&cart.ID, &cart.Status, &cart.CustomerID, &cart.Created, &cart.CreatedBy, &cart.Updated, &cart.UpdatedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, response.Error("failed to scan cart", err)
	}
	if updated.Valid {
		base.Updated.Time = updated.Time
	}
	if !base.UpdatedBy.Valid {
		base.UpdatedBy.String = ""
	}

	cart.Base = base

	return cart, nil
}
