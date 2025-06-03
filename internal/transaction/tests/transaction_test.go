package test

import (
	customer "Dzaakk/simple-commerce/internal/customer/repository"
	product "Dzaakk/simple-commerce/internal/product/repository"
	scModel "Dzaakk/simple-commerce/internal/shopping_cart/models"
	shoppingCart "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	sc "Dzaakk/simple-commerce/internal/shopping_cart/usecase"
	"Dzaakk/simple-commerce/internal/transaction/models"
	repo "Dzaakk/simple-commerce/internal/transaction/repository"
	usecase "Dzaakk/simple-commerce/internal/transaction/usecase"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var psql *sql.DB

func TestDoubleTransaction(t *testing.T) {

	repoCartItem := shoppingCart.NewShoppingCartItemRepository(psql)
	repoCustomer := customer.NewCustomerRepository(psql)
	repoProduct := product.NewProductRepository(psql)
	repoCart := shoppingCart.NewShoppingCartRepository(psql)
	repo := repo.NewTransactionRepository(psql)
	var db *sql.DB

	usecase := usecase.NewTransactionUseCase(repo, repoCart, repoCartItem, repoCustomer, repoProduct, db)

	_ = repoProduct.SetStockById(10, 1) // set stock product with id 10 to 1

	cartUsecase := sc.NewShoppingCartUseCase(repoCart, repoCartItem, repoProduct)

	customer1CartReq := scModel.ShoppingCartReq{
		CustomerId: "2",
		ProductId:  "10",
		Quantity:   "1",
	}
	customer2CartReq := scModel.ShoppingCartReq{
		CustomerId: "3",
		ProductId:  "10",
		Quantity:   "1",
	}

	cart1, _ := cartUsecase.Add(customer1CartReq)
	cart2, _ := cartUsecase.Add(customer2CartReq)

	transaction1 := models.TransactionReq{
		CustomerId: "1",
		CartId:     cart1.NewCartId,
	}
	transaction2 := models.TransactionReq{
		CustomerId: "2",
		CartId:     cart2.NewCartId,
	}

	var wg sync.WaitGroup
	var err1, err2 error

	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := usecase.CreateTransaction(transaction1)
		if err != nil {
			err1 = err
		}
	}()
	go func() {
		defer wg.Done()
		_, err := usecase.CreateTransaction(transaction2)
		if err != nil {
			err2 = err
		}
	}()

	wg.Wait()
	cartId1, _ := strconv.Atoi(cart1.NewCartId)
	cartId2, _ := strconv.Atoi(cart2.NewCartId)
	// Check the results
	if err1 == nil && err2 != nil {
		log.Print("Transaction Success on customer 1")

		_ = repoCartItem.DeleteAll(cartId2)
		_ = repoCart.DeleteShoppingCart(cartId2)

		assert.NoError(t, err1)
		assert.Error(t, err2, "Expected second transaction to fail due to insufficient stock")
	} else if err1 != nil && err2 == nil {
		log.Print("Transaction Success on customer 2")

		_ = repoCartItem.DeleteAll(cartId1)
		_ = repoCart.DeleteShoppingCart(cartId1)

		assert.NoError(t, err2)
		assert.Error(t, err1, "Expected first transaction to fail due to insufficient stock")
	} else {
		t.Error("Both transactions should not succeed or fail")
	}

	_ = repoProduct.SetStockById(10, 1) // set back product with id 10 to 1
}

func TestRaceTransaction(t *testing.T) {
	for i := 0; i < 50; i++ { // Run TestRaceCondition 50 times
		t.Run(fmt.Sprintf("Run-%d", i+1), func(t *testing.T) {
			TestDoubleTransaction(t)
		})
	}
}
