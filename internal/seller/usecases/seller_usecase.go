package usecases

import model "Dzaakk/simple-commerce/internal/seller/models"

type SellerUseCase interface {
	Create(data model.SellerReq) (int64, error)
	Update(data model.SellerReq) (int64, error)
	FindById(sellerId int) (*model.SellerRes, error)
	FindByUsername(username string) (*model.SellerRes, error)
	Deactivate(sellerId int) (int64, error)
}
