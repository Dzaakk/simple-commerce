package usecase

import "Dzaakk/simple-commerce/internal/transaction/models"

type TransactionItemUseCase interface {
	GetListTransaction(customerId int)
	GetTransactionItemDetail(transactionItemId int64)
	CreateTransactionItemDetail(transactionItem models.TTransactionItem)
}
