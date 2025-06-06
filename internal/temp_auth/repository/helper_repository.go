package repository

import (
	"Dzaakk/simple-commerce/internal/auth/model"
	"database/sql"
	"errors"
)

func rowsToCustomerActivationCode(rows *sql.Rows) (*model.TCustomerActivationCode, error) {
	ac := model.TCustomerActivationCode{}
	err := rows.Scan(&ac.CustomerID, &ac.CodeActivation, &ac.IsUsed, &ac.CreatedAt, &ac.UsedAt)
	if err != nil {
		return nil, err
	}

	return &ac, nil
}

func retrieveCustomerCodeActivaton(rows *sql.Rows) (*model.TCustomerActivationCode, error) {
	if rows.Next() {
		return rowsToCustomerActivationCode(rows)
	}
	return nil, errors.New("code activation not found")
}

func rowsToSellerActivationCode(rows *sql.Rows) (*model.TSellerActivationCode, error) {
	sc := model.TSellerActivationCode{}
	err := rows.Scan(&sc.SellerID, &sc.CodeActivation, &sc.IsUsed, &sc.CreatedAt, &sc.UsedAt)
	if err != nil {
		return nil, err
	}

	return &sc, nil
}

func retrieveSellerCodeActivaton(rows *sql.Rows) (*model.TSellerActivationCode, error) {
	if rows.Next() {
		return rowsToSellerActivationCode(rows)
	}
	return nil, errors.New("code activation not found")
}
