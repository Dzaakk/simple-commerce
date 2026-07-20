package service

import (
	"Dzaakk/simple-commerce/internal/catalog/model"
	"Dzaakk/simple-commerce/package/response"
	"context"
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

func (s *InventoryServiceImpl) ReserveStock(ctx context.Context, productID string, qty int) error {
	if productID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter product id")
	}
	if qty <= 0 {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter quantity")
	}

	return s.repo.ReserveStock(ctx, productID, qty)
}

func (s *InventoryServiceImpl) ReleaseStock(ctx context.Context, productID string, qty int) error {
	if productID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter product id")
	}
	if qty <= 0 {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter quantity")
	}

	return s.repo.ReleaseStock(ctx, productID, qty)
}
