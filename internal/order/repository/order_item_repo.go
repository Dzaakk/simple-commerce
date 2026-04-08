package repository

import (
	"Dzaakk/simple-commerce/internal/order/model"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

const (
	orderItemSelectColumns = "id, order_id, product_id, seller_id, quantity, price, subtotal, created_at"

	orderItemQueryFindByOrderID = "SELECT " + orderItemSelectColumns + " FROM public.order_items WHERE order_id=$1 ORDER BY id ASC"
)

type OrderItemRepository struct {
	DB *sql.DB
}

func NewOrderItemRepository(db *sql.DB) *OrderItemRepository {
	return &OrderItemRepository{DB: db}
}

func (r *OrderItemRepository) CreateBatch(ctx context.Context, tx *sql.Tx, items []*model.OrderItem) error {
	if tx == nil {
		return errors.New("transaction is required")
	}
	if len(items) == 0 {
		return nil
	}

	values := make([]string, 0, len(items))
	args := make([]any, 0, len(items)*7)
	argPos := 1

	for _, item := range items {
		values = append(values, "("+placeholders(argPos, 7)+")")
		args = append(args,
			item.OrderID,
			item.ProductID,
			item.SellerID,
			item.Quantity,
			item.Price,
			item.Subtotal,
			item.CreatedAt,
		)
		argPos += 7
	}

	query := formatOrderItemInsert(values)
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return response.ExecError("create order items", err)
	}

	return nil
}

func (r *OrderItemRepository) FindByOrderID(ctx context.Context, orderID string) ([]*model.OrderItem, error) {
	rows, err := r.DB.QueryContext(ctx, orderItemQueryFindByOrderID, orderID)
	if err != nil {
		return nil, response.Error("failed to query order items", err)
	}
	defer rows.Close()

	var items []*model.OrderItem

	for rows.Next() {
		var item model.OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.SellerID,
			&item.Quantity,
			&item.Price,
			&item.Subtotal,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, response.Error("failed to scan order item", err)
		}

		items = append(items, &item)
	}

	return items, nil
}

func placeholders(start, count int) string {
	parts := make([]string, 0, count)
	for i := 0; i < count; i++ {
		parts = append(parts, "$"+strconv.Itoa(start+i))
	}
	return strings.Join(parts, ", ")
}

func formatOrderItemInsert(values []string) string {
	return "INSERT INTO public.order_items (order_id, product_id, seller_id, quantity, price, subtotal, created_at) VALUES " + strings.Join(values, ", ")
}
