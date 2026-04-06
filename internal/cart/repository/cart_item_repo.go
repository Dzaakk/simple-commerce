package repository

import (
	"Dzaakk/simple-commerce/internal/cart/model"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"time"
)

const (
	cartItemSelectColumns = "id, cart_id, product_id, quantity, price_snapshot, created_at, updated_at"

	cartItemQueryFindByCartID = "SELECT " + cartItemSelectColumns + " FROM public.cart_items WHERE cart_id=$1 ORDER BY created_at ASC"
	cartItemQueryUpsert       = "INSERT INTO public.cart_items (cart_id, product_id, quantity, price_snapshot, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (cart_id, product_id) DO UPDATE SET quantity=EXCLUDED.quantity, price_snapshot=EXCLUDED.price_snapshot, updated_at=EXCLUDED.updated_at"
	cartItemQueryDelete       = "DELETE FROM public.cart_items WHERE cart_id=$1 AND product_id=$2"
	cartItemQueryClear        = "DELETE FROM public.cart_items WHERE cart_id=$1"
)

type CartItemRepository struct {
	DB *sql.DB
}

func NewCartItemRepository(db *sql.DB) *CartItemRepository {
	return &CartItemRepository{DB: db}
}

func (r *CartItemRepository) GetCartItems(ctx context.Context, cartID string) ([]*model.CartItem, error) {
	rows, err := r.DB.QueryContext(ctx, cartItemQueryFindByCartID, cartID)
	if err != nil {
		return nil, response.Error("failed to query cart items", err)
	}
	defer rows.Close()

	var items []*model.CartItem

	for rows.Next() {
		var item model.CartItem
		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
			&item.PriceSnapshot,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, response.Error("failed to scan cart item", err)
		}

		items = append(items, &item)
	}

	return items, nil
}

func (r *CartItemRepository) UpsertItem(ctx context.Context, cartID string, productID string, quantity int, priceSnapshot float64) error {
	result, err := r.DB.ExecContext(
		ctx,
		cartItemQueryUpsert,
		cartID,
		productID,
		quantity,
		priceSnapshot,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return response.ExecError("upsert cart item", err)
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

func (r *CartItemRepository) DeleteItem(ctx context.Context, cartID string, productID string) error {
	result, err := r.DB.ExecContext(ctx, cartItemQueryDelete, cartID, productID)
	if err != nil {
		return response.ExecError("delete cart item", err)
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

func (r *CartItemRepository) ClearItems(ctx context.Context, cartID string) error {
	_, err := r.DB.ExecContext(ctx, cartItemQueryClear, cartID)
	if err != nil {
		return response.ExecError("clear cart items", err)
	}

	return nil
}
