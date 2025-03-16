package usecases

import (
	model "Dzaakk/simple-commerce/internal/seller/models"
	"context"
)

type SellerUseCase interface {
	Create(ctx context.Context, data model.ReqCreate) (int64, error)
	Update(ctx context.Context, data model.ReqUpdate) (int64, error)
	FindById(ctx context.Context, sellerId int64) (*model.ResData, error)
	FindByUsername(ctx context.Context, username string) (*model.ResData, error)
	FindByEmail(ctx context.Context, email string) (*model.TSeller, error)
	Deactivate(ctx context.Context, sellerId int64) (int64, error)
	ChangePassword(ctx context.Context, sellerId int64, newPassword string) (int64, error)
}
