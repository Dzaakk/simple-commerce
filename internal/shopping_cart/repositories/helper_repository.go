package repositories

import (
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	response "Dzaakk/simple-commerce/package/response"
	template "Dzaakk/simple-commerce/package/templates"
	"database/sql"
	"errors"
)

func scanCart(row *sql.Row) (*model.TShoppingCart, error) {
	cart := &model.TShoppingCart{}
	base := template.Base{}
	var updated sql.NullTime

	err := row.Scan(
		&cart.Id, &cart.Status, &cart.CustomerId, &cart.Created, &cart.CreatedBy, &cart.Updated, &cart.UpdatedBy)
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
