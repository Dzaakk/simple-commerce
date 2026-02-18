package repository

import (
	"Dzaakk/simple-commerce/internal/user/domain"
	response "Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
)

const (
	customerSelectColumns = "id, email, password_hash, full_name, phone, status, created_at, updated_at"
	queryCreate           = "INSERT INTO public.customers (id, email, password_hash, full_name, phone, status, created_at, updated_at) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7) RETURNING id"
	queryFindByEmail      = "SELECT " + customerSelectColumns + " FROM public.customers WHERE email=$1"
	queryFindByID         = "SELECT " + customerSelectColumns + " FROM public.customers WHERE id=$1"
	queryUpdate           = "UPDATE public.customers SET email=$1, full_name=$2, phone=$3, status=$4, updated_at=$5 WHERE id=$6"
)

type CustomerRepository struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{DB: db}
}

func (r *CustomerRepository) Create(ctx context.Context, data *domain.Customer) (string, error) {
	var id string

	err := r.DB.QueryRowContext(
		ctx,
		queryCreate,
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

func (r *CustomerRepository) Update(ctx context.Context, data *domain.Customer) (int64, error) {
	result, err := r.DB.ExecContext(
		ctx, queryUpdate,
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

func (r *CustomerRepository) FindByID(ctx context.Context, customerID string) (*domain.Customer, error) {
	row := r.DB.QueryRowContext(ctx, queryFindByID, customerID)

	return scanCustomer(row)
}

func (r *CustomerRepository) FindByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	row := r.DB.QueryRowContext(ctx, queryFindByEmail, email)

	return scanCustomer(row)
}
