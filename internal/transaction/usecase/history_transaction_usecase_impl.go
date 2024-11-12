package usecase

import (
	"Dzaakk/simple-commerce/internal/transaction/models"
	repo "Dzaakk/simple-commerce/internal/transaction/repository"
)

type HistoryTransactionUseCaseImpl struct {
	repo repo.HistoryTransactionRepository
}

func NewHistoryTransactionUseCase(repo repo.HistoryTransactionRepository) HistoryTransactionUseCase {
	return &HistoryTransactionUseCaseImpl{repo}
}

func (t *HistoryTransactionUseCaseImpl) CreateTransactionItemDetail(transactionItem models.THistoryTransaction) {
	panic("unimplemented")
}

func (t *HistoryTransactionUseCaseImpl) GetListTransaction(customerId int) {
	panic("unimplemented")
}

func (t *HistoryTransactionUseCaseImpl) GetTransactionItemDetail(transactionItemId int64) {
	panic("unimplemented")
}
