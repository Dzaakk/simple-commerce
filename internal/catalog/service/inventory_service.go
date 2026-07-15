package service

import (
	"Dzaakk/simple-commerce/internal/catalog/model"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"database/sql"
	"net/http"
)

type InventoryServiceImpl struct {
	repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) *InventoryServiceImpl {
	return &InventoryServiceImpl{repo: repo}
}

func (s *InventoryServiceImpl) FindByProductID(ctx context.Context, productID string) (*model.Inventory, error) {
	if productID == "" {
		return nil, response.NewAppError(http.StatusBadRequest, "invalid parameter product id")
	}

	return s.repo.FindByProductID(ctx, productID)
}

func (s *InventoryServiceImpl) ReserveStock(ctx context.Context, tx *sql.Tx, productID string, qty int) error {
	if productID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter product id")
	}
	if qty <= 0 {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter quantity")
	}
	if tx == nil {
		return response.NewAppError(http.StatusInternalServerError, "internal server error")
	}

	return s.repo.ReserveStock(ctx, tx, productID, qty)
}

func (s *InventoryServiceImpl) ReleaseStock(ctx context.Context, tx *sql.Tx, productID string, qty int) error {
	if productID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter product id")
	}
	if qty <= 0 {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter quantity")
	}
	if tx == nil {
		return response.NewAppError(http.StatusInternalServerError, "internal server error")
	}

	return s.repo.ReleaseStock(ctx, tx, productID, qty)
}
