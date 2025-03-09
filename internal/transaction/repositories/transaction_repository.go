package repositories

import (
	modelItem "Dzaakk/simple-commerce/internal/shopping_cart/models"
	model "Dzaakk/simple-commerce/internal/transaction/models"
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
	FindByCustomerID(ctx context.Context, customerId int64) ([]*model.THistoryTransaction, error)
}
