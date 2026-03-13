package repository

import (
	"Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/constant"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
)

const (
	customerSelectColumns     = "id, email, password_hash, full_name, phone, status, created_at, updated_at"
	customerQueryCreate       = "INSERT INTO public.customers (id, email, password_hash, full_name, phone, status, created_at, updated_at) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7) RETURNING id"
	customerQueryFindByEmail  = "SELECT " + customerSelectColumns + " FROM public.customers WHERE email=$1"
	customerQueryFindByID     = "SELECT " + customerSelectColumns + " FROM public.customers WHERE id=$1"
	customerQueryUpdate       = "UPDATE public.customers SET email=$1, full_name=$2, phone=$3, status=$4, updated_at=$5 WHERE id=$6"
	customerQueryUpdateStatus = "UPDATE public.customers SET status=$1, updated_at=NOW() WHERE id=$2"
)

type CustomerRepository struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{DB: db}
}

func (r *CustomerRepository) Create(ctx context.Context, data *model.Customer) (string, error) {
	var id string

	err := r.DB.QueryRowContext(
		ctx,
		customerQueryCreate,
		data.Email,
		data.PasswordHash,
		data.FullName,
		data.Phone,
		data.Status,
		data.CreatedAt,
		data.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return "", response.Error("failed to create customer", err)
	}

	return id, nil
}

func (r *CustomerRepository) Update(ctx context.Context, data *model.Customer) (int64, error) {
	result, err := r.DB.ExecContext(
		ctx, customerQueryUpdate,
		data.Email,
		data.FullName,
		data.Phone,
		data.Status,
		data.UpdatedAt,
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

func (r *CustomerRepository) FindByID(ctx context.Context, customerID string) (*model.Customer, error) {
	row := r.DB.QueryRowContext(ctx, customerQueryFindByID, customerID)

	return scanCustomer(row)
}

func (r *CustomerRepository) FindByEmail(ctx context.Context, email string) (*model.Customer, error) {
	row := r.DB.QueryRowContext(ctx, customerQueryFindByEmail, email)

	return scanCustomer(row)
}

func (r *CustomerRepository) UpdateStatus(ctx context.Context, customerID string, status constant.UserStatus) error {
	result, err := r.DB.ExecContext(ctx, customerQueryUpdateStatus, status, customerID)
	if err != nil {
		return response.ExecError("update customer status", err)
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
