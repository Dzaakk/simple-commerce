package usecase

type OrderUseCase interface {
	CreateOrder()
	CancelOrder()
	GetListCustomerOrder()
	GetDetailCustomerOrder()
}
