package transaction

import (
	model "Dzaakk/simple-commerce/internal/transaction/models"
)

type TransactionRepository interface {
	Create(data model.TTransaction) (*model.TTransaction, error)
}
