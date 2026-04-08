package service

import (
	"Dzaakk/simple-commerce/internal/catalog/model"
	"context"
	"database/sql"
	"errors"
)

type InventoryServiceImpl struct {
	Repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) InventoryService {
	return &InventoryServiceImpl{Repo: repo}
}

func (s *InventoryServiceImpl) FindByProductID(ctx context.Context, productID string) (*model.Inventory, error) {
	if productID == "" {
		return nil, errors.New("invalid parameter product id")
	}

	return s.Repo.FindByProductID(ctx, productID)
}

func (s *InventoryServiceImpl) ReserveStock(ctx context.Context, tx *sql.Tx, productID string, qty int) error {
	if productID == "" {
		return errors.New("invalid parameter product id")
	}
	if qty <= 0 {
		return errors.New("invalid parameter quantity")
	}
	if tx == nil {
		return errors.New("transaction is required")
	}

	return s.Repo.ReserveStock(ctx, tx, productID, qty)
}

func (s *InventoryServiceImpl) ReleaseStock(ctx context.Context, tx *sql.Tx, productID string, qty int) error {
	if productID == "" {
		return errors.New("invalid parameter product id")
	}
	if qty <= 0 {
		return errors.New("invalid parameter quantity")
	}
	if tx == nil {
		return errors.New("transaction is required")
	}

	return s.Repo.ReleaseStock(ctx, tx, productID, qty)
}
