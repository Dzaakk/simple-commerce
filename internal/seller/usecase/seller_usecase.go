package usecase

import (
	"Dzaakk/simple-commerce/internal/seller/model"
	"context"
)

type SellerUseCase interface {
	Create(ctx context.Context, data model.ReqCreate) (int64, error)
	Update(ctx context.Context, data model.ReqUpdate) (int64, error)
	FindAll(ctx context.Context) ([]*model.ResData, error)
	FindBySellerID(ctx context.Context, sellerID int64) (*model.ResData, error)
	FindByUsername(ctx context.Context, username string) (*model.ResData, error)
	FindByEmail(ctx context.Context, email string) (*model.TSeller, error)
	Deactivate(ctx context.Context, sellerID int64) (int64, error)
	ChangePassword(ctx context.Context, sellerID int64, newPassword string) (int64, error)
}
