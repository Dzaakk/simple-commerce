package usecase

import (
	customer "Dzaakk/simple-commerce/internal/customer/repositories"
	product "Dzaakk/simple-commerce/internal/product/repositories"
	modelItem "Dzaakk/simple-commerce/internal/shopping_cart/models"
	shoppingCart "Dzaakk/simple-commerce/internal/shopping_cart/repositories"
	model "Dzaakk/simple-commerce/internal/transaction/models"
	repo "Dzaakk/simple-commerce/internal/transaction/repositories"
	"Dzaakk/simple-commerce/package/template"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
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

func (t *TransactionUseCaseImpl) CreateTransaction(data model.TransactionReq) (*model.TransactionRes, error) {
	tx, err := t.repo.BeginTransaction()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	cartId, err := strconv.Atoi(data.CartId)
	if err != nil {
		return nil, fmt.Errorf("invalid data: %v", err)
	}
	customerId, err := strconv.Atoi(data.CustomerId)
	if err != nil {
		return nil, fmt.Errorf("invalid data: %v", err)
	}

	listItem, err := t.repoCartItem.RetrieveCartItemsByCartIdWithTx(tx, cartId) // get all items on cart
	if err != nil {
		return nil, err
	}

	res, err := generateReceipt(listItem) // generate receipt and calculate total transaction
	if err != nil {
		return nil, err
	}

	customer, err := t.repoCustomer.GetBalanceWithTx(tx, customerId) // check customer current balance with locking
	if err != nil {
		return nil, err
	}
	totalTransaction, _ := strconv.Atoi(res.TotalTransaction)
	if totalTransaction > int(customer.Balance) {
		return nil, errors.New("insufficient balance")
	}

	err = t.repoCart.UpdateStatusByIdWithTx(tx, cartId, "Paid", data.CustomerId) // update cart status to 'Paid'
	if err != nil {
		return nil, err
	}

	emptyProducts, err := t.repoProduct.UpdateStockWithTx(tx, listItem) // update stock and get list fo empty product
	if err != nil {
		return nil, err
	}

	err = t.repoCartItem.SetQuantityWithTx(tx, emptyProducts)
	if err != nil {
		return nil, err
	}

	newBalance := customer.Balance - float64(totalTransaction)
	err = t.repoCustomer.UpdateBalanceWithTx(tx, customerId, newBalance) // update balance customer
	if err != nil {
		return nil, err
	}

	transactionDate, err := insertToTableTransactionWithTx(tx, t, customerId, cartId, totalTransaction) // insert to table transaction
	if err != nil {
		return nil, err
	}

	err = t.repoCartItem.DeleteAllWithTx(tx, cartId) // delete all item on cart item
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	res.CustomerId = data.CustomerId
	res.TransactionDate = *transactionDate
	return res, nil
}

func insertToTableTransactionWithTx(tx *sql.Tx, t *TransactionUseCaseImpl, customerId, cartId, totalTransaction int) (*string, error) {
	newTransaction := model.TTransaction{
		CustomerId:      customerId,
		CartId:          cartId,
		TotalAmount:     float32(totalTransaction),
		TransactionDate: time.Now(),
		Status:          "Success",
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: fmt.Sprintf("%d", customerId),
		},
	}

	data, err := t.repo.CreateWithTx(tx, newTransaction)
	if err != nil {
		return nil, err
	}
	transactionDate := data.TransactionDate.Format("06-01-02 15:04:05")

	return &transactionDate, nil
}

func generateReceipt(listItem []*modelItem.TCartItemDetail) (*model.TransactionRes, error) {
	var res model.TransactionRes
	var listProduct []model.ListProduct
	total := 0
	for _, item := range listItem {
		product := model.ListProduct{
			ProductName: item.ProductName,
			Price:       fmt.Sprintf("%.0f", item.Price),
			Quantity:    fmt.Sprintf("%d", item.Quantity),
		}
		listProduct = append(listProduct, product)
		total = total + (int(item.Price) * item.Quantity)
	}
	res.ListProduct = listProduct
	res.TotalTransaction = fmt.Sprintf("%d", total)

	return &res, nil
}

func (t *TransactionUseCaseImpl) GetTransaction(customerId int64) ([]*model.CustomerListTransactionRes, error) {
	panic("unimplemented")
}

// GetDetailTransaction implements TransactionUseCase.
func (t *TransactionUseCaseImpl) GetDetailTransaction(transactionId int64) ([]*model.CustomerListTransactionRes, error) {
	panic("unimplemented")
}
