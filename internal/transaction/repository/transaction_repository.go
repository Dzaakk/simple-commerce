package repository

import (
	modelItem "Dzaakk/simple-commerce/internal/shopping_cart/model"
	"Dzaakk/simple-commerce/internal/transaction/model"
	"context"
	"database/sql"
)

type TransactionRepository interface {
	Create(ctx context.Context, data model.TTransaction) (*model.TTransaction, error)
	CreateWithTx(ctx context.Context, tx *sql.Tx, data model.TTransaction) (*model.TTransaction, error)
	BeginTransaction() (*sql.Tx, error)
}

type HistoryTransactionRepository interface {
	Create(ctx context.Context, data []*modelItem.TCartItemDetail, customerId int64) error
	FindByCustomerID(ctx context.Context, customerID int64) ([]*model.THistoryTransaction, error)
}
