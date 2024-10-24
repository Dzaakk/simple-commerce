package repository

import (
	model "Dzaakk/simple-commerce/internal/transaction/models"
)

type TransactionRepository interface {
	Create(data model.TTransaction) (*model.TTransaction, error)
	FindByCustomerId(customerId int64)
}
