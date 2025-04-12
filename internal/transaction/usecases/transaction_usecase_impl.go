package usecases

import (
	customer "Dzaakk/simple-commerce/internal/customer/repositories"
	product "Dzaakk/simple-commerce/internal/product/repositories"
	shoppingCart "Dzaakk/simple-commerce/internal/shopping_cart/repositories"
	model "Dzaakk/simple-commerce/internal/transaction/models"
	repo "Dzaakk/simple-commerce/internal/transaction/repositories"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

type TransactionUseCaseImpl struct {
	db           *sql.DB
	repo         repo.TransactionRepository
	repoCart     shoppingCart.ShoppingCartRepository
	repoCartItem shoppingCart.ShoppingCartItemRepository
	repoCustomer customer.CustomerRepository
	repoProduct  product.ProductRepository
}

func NewTransactionUseCase(repo repo.TransactionRepository, repoCart shoppingCart.ShoppingCartRepository, repoCartItem shoppingCart.ShoppingCartItemRepository, repoCustomer customer.CustomerRepository, repoProduct product.ProductRepository, db *sql.DB) TransactionUseCase {

	return &TransactionUseCaseImpl{db, repo, repoCart, repoCartItem, repoCustomer, repoProduct}
}

func (t *TransactionUseCaseImpl) CreateTransaction(ctx context.Context, data model.TransactionReq) (*model.TransactionRes, error) {
	tx, err := t.repo.BeginTransaction()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	cartID, err := strconv.Atoi(data.CartID)
	if err != nil {
		return nil, fmt.Errorf("invalid data: %v", err)
	}
	customerID, err := strconv.Atoi(data.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("invalid data: %v", err)
	}

	listItem, err := t.repoCartItem.RetrieveCartItemsByCartIDWithTx(ctx, tx, cartID) // get all items on cart
	if err != nil {
		return nil, err
	}

	res, err := generateReceipt(listItem) // generate receipt and calculate total transaction
	if err != nil {
		return nil, err
	}

	customer, err := t.repoCustomer.GetBalanceWithTx(ctx, tx, int64(customerID)) // check customer current balance with locking
	if err != nil {
		return nil, err
	}
	totalTransaction, _ := strconv.Atoi(res.TotalTransaction)
	if totalTransaction > int(customer.Balance) {
		return nil, errors.New("insufficient balance")
	}

	err = t.repoCart.UpdateStatusByCartIDWithTx(ctx, tx, cartID, "Paid", data.CustomerID) // update cart status to 'Paid'
	if err != nil {
		return nil, err
	}

	emptyProducts, err := t.repoProduct.UpdateStockWithTx(ctx, tx, listItem) // update stock and get list fo empty product
	if err != nil {
		return nil, err
	}

	err = t.repoCartItem.SetEmptyQuantityWithTx(ctx, tx, emptyProducts)
	if err != nil {
		return nil, err
	}

	newBalance := customer.Balance - float64(totalTransaction)
	err = t.repoCustomer.UpdateBalanceWithTx(ctx, tx, int64(customerID), newBalance) // update balance customer
	if err != nil {
		return nil, err
	}

	transactionDate, err := insertToTableTransactionWithTx(ctx, tx, t, customerID, cartID, totalTransaction) // insert to table transaction
	if err != nil {
		return nil, err
	}

	err = t.repoCartItem.DeleteAllWithTx(ctx, tx, cartID) // delete all item on cart item
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	res.CustomerID = data.CustomerID
	res.TransactionDate = *transactionDate
	return res, nil
}

func (t *TransactionUseCaseImpl) GetTransaction(ctx context.Context, customerID int64) ([]*model.CustomerListTransactionRes, error) {
	panic("unimplemented")
}

// GetDetailTransaction implements TransactionUseCase.
func (t *TransactionUseCaseImpl) GetDetailTransaction(ctx context.Context, transactionID int64) ([]*model.CustomerListTransactionRes, error) {
	panic("unimplemented")
}
