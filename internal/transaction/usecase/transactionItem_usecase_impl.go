package usecase

import (
	"Dzaakk/simple-commerce/internal/transaction/models"
	repo "Dzaakk/simple-commerce/internal/transaction/repository"
)

type TransactionItemUseCaseImpl struct {
	repo repo.TransactionItemRepository
}

func NewTransactionItemUseCase(repo repo.TransactionItemRepository) TransactionItemUseCase {
	return &TransactionItemUseCaseImpl{repo}
}

// CreateTransactionItemDetail implements TransactionItemUseCase.
func (t *TransactionItemUseCaseImpl) CreateTransactionItemDetail(transactionItem models.TTransactionItem) {
	panic("unimplemented")
}

// GetListTransaction implements TransactionItemUseCase.
func (t *TransactionItemUseCaseImpl) GetListTransaction(customerId int) {
	panic("unimplemented")
}

// GetTransactionItemDetail implements TransactionItemUseCase.
func (t *TransactionItemUseCaseImpl) GetTransactionItemDetail(transactionItemId int64) {
	panic("unimplemented")
}
