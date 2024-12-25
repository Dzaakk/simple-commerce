package usecases

import model "Dzaakk/simple-commerce/internal/seller/models"

type SellerUseCase interface {
	Create(data model.ReqCreate) (int64, error)
	Update(data model.ReqUpdate) (int64, error)
	FindById(sellerId int) (*model.ResData, error)
	FindByUsername(username string) (*model.ResData, error)
	Deactivate(sellerId int) (int64, error)
}
