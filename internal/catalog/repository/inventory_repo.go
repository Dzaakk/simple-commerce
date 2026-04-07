package repository

import (
	catalogModel "Dzaakk/simple-commerce/internal/catalog/model"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"errors"
)

const (
	inventorySelectColumns     = "id, product_id, stock_quantity, reserved_quantity, version, created_at, updated_at"
	inventoryQueryFindByProduct = "SELECT " + inventorySelectColumns + " FROM public.inventories WHERE product_id=$1"
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
