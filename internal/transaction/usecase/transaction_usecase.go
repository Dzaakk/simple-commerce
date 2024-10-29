package usecase

import (
	model "Dzaakk/simple-commerce/internal/transaction/models"
)

type TransactionUseCase interface {
	CreateTransaction(data model.TransactionReq) (*model.TransactionRes, error)
	GetTransaction(customerId int64) ([]*model.CustomerListTransactionRes, error)
	GetDetailTransaction(transactionId int64) ([]*model.CustomerListTransactionRes, error)
}
type TransactionItemUseCase interface {
	GetListTransaction(customerId int)
	GetTransactionItemDetail(transactionItemId int64)
	CreateTransactionItemDetail(transactionItem model.TTransactionItem)
}
