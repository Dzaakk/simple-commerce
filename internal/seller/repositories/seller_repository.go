package repositories

import (
	model "Dzaakk/simple-commerce/internal/seller/models"
	"context"
)

type SellerRepository interface {
	Create(ctx context.Context, data model.TSeller) (int64, error)
	FindById(ctx context.Context, sellerId int64) (*model.TSeller, error)
	FindAll(ctx context.Context) ([]*model.TSeller, error)
	FindByUsername(ctx context.Context, username string) (*model.TSeller, error)
	FindByEmail(ctx context.Context, email string) (*model.TSeller, error)
	Deactive(ctx context.Context, sellerId int64) (int64, error)
	Update(ctx context.Context, data model.TSeller) (int64, error)
	UpdatePassword(ctx context.Context, sellerId int64, newPassword string) (int64, error)
	InsertBalance(ctx context.Context, sellerId, balance int64) error
}
