package transaction

import model "Dzaakk/simple-commerce/internal/transaction/models"

type TransactionUseCase interface {
	CreateTransaction(data model.TransactionReq) (*model.TransactionRes, error)
}
