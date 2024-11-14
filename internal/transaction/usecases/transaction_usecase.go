package usecase

import (
	model "Dzaakk/simple-commerce/internal/transaction/models"
)

type TransactionUseCase interface {
	CreateTransaction(data model.TransactionReq) (*model.TransactionRes, error)
	GetTransaction(customerId int64) ([]*model.CustomerListTransactionRes, error)
	GetDetailTransaction(transactionId int64) ([]*model.CustomerListTransactionRes, error)
}
type HistoryTransactionUseCase interface {
	GetListHistoryTransaction(customerId int64) ([]*model.HistoryTransaction, error)
	GetHistoryTransactionDetail(transactionId int64)
	CreateHistoryTransaction(transactionItem model.THistoryTransaction)
}
