package usecase

import (
	model "Dzaakk/simple-commerce/internal/transaction/models"
	"context"
)

type TransactionUseCase interface {
	CreateTransaction(ctx context.Context, data model.TransactionReq) (*model.TransactionRes, error)
	GetTransaction(ctx context.Context, customerId int64) ([]*model.CustomerListTransactionRes, error)
	GetDetailTransaction(ctx context.Context, transactionId int64) ([]*model.CustomerListTransactionRes, error)
}
type HistoryTransactionUseCase interface {
	GetListHistoryTransaction(ctx context.Context, customerId int64) ([]*model.HistoryTransaction, error)
	GetHistoryTransactionDetail(ctx context.Context, transactionId int64)
	CreateHistoryTransaction(ctx context.Context, transactionItem model.THistoryTransaction)
}
