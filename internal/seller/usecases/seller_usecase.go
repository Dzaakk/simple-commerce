package usecases

type SellerUseCase interface {
	Create()
	Update()
	FindById()
	FindByName()
	Deactivate()
}
