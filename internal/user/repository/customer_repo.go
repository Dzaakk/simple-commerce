package repository

import (
	"Dzaakk/simple-commerce/internal/customer/model"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"errors"
)

const (
	customerSelectColumns = "id, username, email, password, gender, phone_number, balance, status, date_of_birth, profile_picture, last_login, created, created_by, updated, updated_by"
	queryCreate           = "INSERT INTO public.customer (username, email, password, gender, phone_number, balance, status, date_of_birth, profile_picture, last_login, created, created_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id"
	queryFindByEmail      = "SELECT " + customerSelectColumns + " FROM public.customer WHERE email=$1"
	queryFindByID         = "SELECT " + customerSelectColumns + " FROM public.customer WHERE id=$1"
	queryFindBalanceForTx = "SELECT id, balance FROM public.customer WHERE id=$1 FOR UPDATE"
	queryUpdate           = "UPDATE public.customer SET username=$1, email=$2, phone_number=$3, date_of_birth=$4, address=$5, updated_by=$6, updated=NOW() WHERE id=$7"
	queryUpdateBalanceTx  = "UPDATE public.customer SET balance=$1, updated=NOW() WHERE id=$2"
)

type CustomerRepositoryImpl struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepositoryImpl {
	return &CustomerRepositoryImpl{DB: db}
}

func (repo *CustomerRepositoryImpl) Create(ctx context.Context, data *model.Customers) (int64, error) {
	var id int64
	err := repo.DB.QueryRowContext(
		ctx,
		queryCreate,
		data.Username, data.Email, data.Password,
		data.Gender, data.PhoneNumber, data.Balance,
		data.Status, data.DateOfBirth, data.ProfilePicture,
		data.LastLogin, data.Base.Created, data.Base.CreatedBy,
	).Scan(&id)
	if err != nil {
		return 0, response.Error("failed to create customer", err)
	}

	return id, nil
}

func (repo *CustomerRepositoryImpl) Update(ctx context.Context, data *model.Customers) (int64, error) {
	result, err := repo.DB.ExecContext(
		ctx, queryUpdate,
		data.Username, data.Email, data.PhoneNumber,
		data.DateOfBirth, data.Address, data.UpdatedBy,
		data.ID,
	)

	if err != nil {
		return 0, response.ExecError("update customer", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, response.Error("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return 0, response.Error("no rows updated", sql.ErrNoRows)
	}

	return rowsAffected, nil
}

func (repo *CustomerRepositoryImpl) FindByID(ctx context.Context, customerID int64) (*model.Customers, error) {
	row := repo.DB.QueryRowContext(ctx, queryFindByID, customerID)

	return scanCustomer(row)
}

func (repo *CustomerRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.Customers, error) {
	row := repo.DB.QueryRowContext(ctx, queryFindByEmail, email)

	return scanCustomer(row)
}

func (repo *CustomerRepositoryImpl) GetBalanceWithTx(ctx context.Context, tx *sql.Tx, customerID int64) (*model.Customers, error) {
	row := tx.QueryRowContext(ctx, queryFindBalanceForTx, customerID)

	customer := &model.Customers{}
	if err := row.Scan(&customer.ID, &customer.Balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, response.Error("failed to scan customer balance", err)
	}

	return customer, nil
}

func (repo *CustomerRepositoryImpl) UpdateBalanceWithTx(ctx context.Context, tx *sql.Tx, customerID int64, balance float64) error {
	result, err := tx.ExecContext(ctx, queryUpdateBalanceTx, balance, customerID)
	if err != nil {
		return response.ExecError("update customer balance", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return response.Error("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return response.Error("no rows updated", sql.ErrNoRows)
	}

	return nil
}
