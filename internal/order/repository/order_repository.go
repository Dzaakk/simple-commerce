package repository

type OrderRepository interface {
	Create()
	Update()
	FindByID()
	FindByCustomerID()
}
