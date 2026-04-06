package repository

import (
	"Dzaakk/simple-commerce/internal/cart/model"
	"context"
	"database/sql"
	"time"
)

const (
	cartSelectColumns = "id, customer_id, created_at, updated_at"

	cartQueryFindByCustomerID = "SELECT " + cartSelectColumns + " FROM public.shopping_carts WHERE customer_id=$1"
	cartQueryGetOrCreate      = "INSERT INTO public.shopping_carts (customer_id, created_at, updated_at) VALUES ($1, $2, $3) ON CONFLICT (customer_id) DO UPDATE SET updated_at = EXCLUDED.updated_at RETURNING id, customer_id, created_at, updated_at"
)

type CartRepository struct {
	DB *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{DB: db}
}

func (r *CartRepository) GetCartByCustomerID(ctx context.Context, customerID string) (*model.Cart, error) {
	row := r.DB.QueryRowContext(ctx, cartQueryFindByCustomerID, customerID)

	return scanCart(row)
}

func (r *CartRepository) GetOrCreateCart(ctx context.Context, customerID string) (*model.Cart, error) {
	now := time.Now()
	row := r.DB.QueryRowContext(ctx, cartQueryGetOrCreate, customerID, now, now)

	return scanCart(row)
}
