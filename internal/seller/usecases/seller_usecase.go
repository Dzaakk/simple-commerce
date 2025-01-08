package usecases

import model "Dzaakk/simple-commerce/internal/seller/models"

type SellerUseCase interface {
	Create(data model.ReqCreate) (int64, error)
	Update(data model.ReqUpdate) (int64, error)
	FindById(sellerId int64) (*model.ResData, error)
	FindByUsername(username string) (*model.ResData, error)
	Deactivate(sellerId int64) (int64, error)
	ChangePassword(sellerId int64, newPassword string) (int64, error)
	Login() error
}
