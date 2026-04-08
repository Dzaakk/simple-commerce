package repository

import (
	"Dzaakk/simple-commerce/internal/order/dto"
	"Dzaakk/simple-commerce/internal/order/model"
	"Dzaakk/simple-commerce/package/constant"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	orderSelectColumns = "id, order_number, customer_id, status, total_amount, shipping_address, created_at, updated_at"

	orderQueryCreate      = "INSERT INTO public.orders (order_number, customer_id, status, total_amount, shipping_address, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	orderQueryFindByID    = "SELECT " + orderSelectColumns + " FROM public.orders WHERE id=$1"
	orderQueryUpdateStatus = "UPDATE public.orders SET status=$1, updated_at=$2 WHERE id=$3"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) Create(ctx context.Context, tx *sql.Tx, data *model.Order) (string, error) {
	if tx == nil {
		return "", errors.New("transaction is required")
	}

	var id string
	row := tx.QueryRowContext(
		ctx,
		orderQueryCreate,
		data.OrderNumber,
		data.CustomerID,
		data.Status,
		data.TotalAmount,
		data.ShippingAddress,
		data.CreatedAt,
		data.UpdatedAt,
	)
	if err := row.Scan(&id); err != nil {
		return "", response.Error("failed to create order", err)
	}

	return id, nil
}

func (r *OrderRepository) FindByID(ctx context.Context, orderID string) (*model.Order, error) {
	row := r.DB.QueryRowContext(ctx, orderQueryFindByID, orderID)

	return scanOrder(row)
}

func (r *OrderRepository) FindByCustomerID(ctx context.Context, customerID string, filter dto.OrderFilter) ([]*model.Order, error) {
	query, args := buildOrderQuery(customerID, filter)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, response.Error("failed to query orders", err)
	}
	defer rows.Close()

	var orders []*model.Order

	for rows.Next() {
		var o model.Order
		err := rows.Scan(
			&o.ID,
			&o.OrderNumber,
			&o.CustomerID,
			&o.Status,
			&o.TotalAmount,
			&o.ShippingAddress,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if err != nil {
			return nil, response.Error("failed to scan order", err)
		}

		orders = append(orders, &o)
	}

	return orders, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, tx *sql.Tx, orderID string, status constant.OrderStatus) error {
	if tx == nil {
		return errors.New("transaction is required")
	}

	result, err := tx.ExecContext(ctx, orderQueryUpdateStatus, status, time.Now(), orderID)
	if err != nil {
		return response.ExecError("update order status", err)
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

func (r *OrderRepository) GenerateOrderNumber(ctx context.Context) (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")

	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)

	var count int
	err := r.DB.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM public.orders WHERE created_at >= $1 AND created_at < $2",
		start,
		end,
	).Scan(&count)
	if err != nil {
		return "", response.Error("failed to count orders", err)
	}

	seq := count + 1
	return fmt.Sprintf("ORD-%s-%04d", dateStr, seq), nil
}

func buildOrderQuery(customerID string, filter dto.OrderFilter) (string, []any) {
	query := "SELECT " + orderSelectColumns + " FROM public.orders WHERE customer_id = $1"
	args := []any{customerID}
	argPos := 2

	if filter.Status != nil && *filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, *filter.Status)
		argPos++
	}

	query += " ORDER BY created_at DESC, id DESC"

	limit := filter.Limit
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, limit)
		argPos++
	}

	if filter.Page > 1 && limit > 0 {
		offset := (filter.Page - 1) * limit
		query += fmt.Sprintf(" OFFSET $%d", argPos)
		args = append(args, offset)
		argPos++
	}

	return query, args
}

func scanOrder(row *sql.Row) (*model.Order, error) {
	var o model.Order

	if err := row.Scan(
		&o.ID,
		&o.OrderNumber,
		&o.CustomerID,
		&o.Status,
		&o.TotalAmount,
		&o.ShippingAddress,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, response.Error("failed to scan order", err)
	}

	return &o, nil
}
