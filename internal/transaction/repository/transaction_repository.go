package repository

import (
	modelItem "Dzaakk/simple-commerce/internal/shopping_cart/models"
	model "Dzaakk/simple-commerce/internal/transaction/models"
)

type TransactionRepository interface {
	Create(data model.TTransaction) (*model.TTransaction, error)
	FindByCustomerId(customerId int64)
}

type TransactionItemRepository interface {
	Create(data []*modelItem.TCartItemDetail, customerId int64) error
}
