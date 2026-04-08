package repository

import (
	catalogModel "Dzaakk/simple-commerce/internal/catalog/model"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"errors"
	"time"
)

const (
	inventorySelectColumns      = "id, product_id, stock_quantity, reserved_quantity, version, created_at, updated_at"
	inventoryQueryFindByProduct = "SELECT " + inventorySelectColumns + " FROM public.inventories WHERE product_id=$1"
	inventoryQueryReserveStock  = `
	UPDATE public.inventories
	SET reserved_quantity = reserved_quantity + $1, updated_at = $2, version = version + 1
	WHERE product_id = $3 AND (stock_quantity - reserved_quantity) >= $1
	`
	inventoryQueryReleaseStock = `
	UPDATE public.inventories
	SET reserved_quantity = reserved_quantity - $1, updated_at = $2, version = version + 1
	WHERE product_id = $3 AND reserved_quantity >= $1
	`
)

type InventoryRepository struct {
	DB *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{DB: db}
}

func (r *InventoryRepository) FindByProductID(ctx context.Context, productID string) (*catalogModel.Inventory, error) {
	row := r.DB.QueryRowContext(ctx, inventoryQueryFindByProduct, productID)

	var inv catalogModel.Inventory
	if err := row.Scan(
		&inv.ID,
		&inv.ProductID,
		&inv.StockQuantity,
		&inv.ReservedQuantity,
		&inv.Version,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, response.Error("failed to scan inventory", err)
	}

	return &inv, nil
}

func (r *InventoryRepository) ReserveStock(ctx context.Context, tx *sql.Tx, productID string, qty int) error {
	if tx == nil {
		return errors.New("transaction is required")
	}
	if qty <= 0 {
		return errors.New("invalid parameter quantity")
	}

	result, err := tx.ExecContext(ctx, inventoryQueryReserveStock, qty, time.Now(), productID)
	if err != nil {
		return response.ExecError("reserve stock", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return response.Error("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return response.Error("insufficient stock", sql.ErrNoRows)
	}

	return nil
}

func (r *InventoryRepository) ReleaseStock(ctx context.Context, tx *sql.Tx, productID string, qty int) error {
	if tx == nil {
		return errors.New("transaction is required")
	}
	if qty <= 0 {
		return errors.New("invalid parameter quantity")
	}

	result, err := tx.ExecContext(ctx, inventoryQueryReleaseStock, qty, time.Now(), productID)
	if err != nil {
		return response.ExecError("release stock", err)
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
