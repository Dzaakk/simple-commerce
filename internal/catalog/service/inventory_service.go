package service

import (
	"Dzaakk/simple-commerce/internal/catalog/model"
	"context"
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
