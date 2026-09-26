package repository

import (
	"context"
	"database/sql"
	"errors"

	"Dzaakk/simple-commerce/internal/user/model"
	"Dzaakk/simple-commerce/package/response"

	"github.com/lib/pq"
)

const (
	customerColumns = "id, email, password_hash, full_name, status, created_at, updated_at"
	createCustomer  = "INSERT INTO users (email, password_hash, full_name, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	findByEmail     = "SELECT " + customerColumns + " FROM users WHERE lower(email) = lower($1)"
	findByID        = "SELECT " + customerColumns + " FROM users WHERE id = $1"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(ctx context.Context, customer *model.Customer) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, createCustomer,
		customer.Email, customer.PasswordHash, customer.FullName, customer.Status,
		customer.CreatedAt, customer.UpdatedAt,
	).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return "", response.ErrEmailAlreadyExist
		}
		return "", response.Error("failed to create customer", err)
	}
	return id, nil
}

func (r *CustomerRepository) FindByEmail(ctx context.Context, email string) (*model.Customer, error) {
	return scanCustomer(r.db.QueryRowContext(ctx, findByEmail, email))
}

func (r *CustomerRepository) FindByID(ctx context.Context, id string) (*model.Customer, error) {
	return scanCustomer(r.db.QueryRowContext(ctx, findByID, id))
}

func scanCustomer(row *sql.Row) (*model.Customer, error) {
	var customer model.Customer
	err := row.Scan(
		&customer.ID, &customer.Email, &customer.PasswordHash, &customer.FullName,
		&customer.Status, &customer.CreatedAt, &customer.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, response.Error("failed to find customer", err)
	}
	return &customer, nil
}
