package repository

import (
	modelItem "Dzaakk/simple-commerce/internal/shopping_cart/models"
	model "Dzaakk/simple-commerce/internal/transaction/models"
	"database/sql"
)

type TransactionRepository interface {
	Create(data model.TTransaction) (*model.TTransaction, error)
	CreateWithTx(tx *sql.Tx, data model.TTransaction) (*model.TTransaction, error)
	FindByCustomerId(customerId int64)
	BeginTransaction() (*sql.Tx, error)
}

type HistoryTransactionRepository interface {
	Create(data []*modelItem.TCartItemDetail, customerId int64) error
	FindByCustomerId(customerId int64) ([]*model.THistoryTransaction, error)
}
