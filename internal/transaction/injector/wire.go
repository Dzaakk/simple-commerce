//go:build wireinject
// +build wireinject

package injector

import (
	customerRepo "Dzaakk/simple-commerce/internal/customer/repository"
	productRepo "Dzaakk/simple-commerce/internal/product/repository"
	cartRepo "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"Dzaakk/simple-commerce/internal/transaction/handler"
	"Dzaakk/simple-commerce/internal/transaction/repository"
	"Dzaakk/simple-commerce/internal/transaction/route"
	"Dzaakk/simple-commerce/internal/transaction/usecase"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *route.TransactionRoutes {
	wire.Build(
		repository.NewTransactionRepository,
		usecase.NewTransactionUseCase,
		handler.NewTransactionHandler,
		route.NewTransactionRoutes,
		cartRepo.NewShoppingCartRepository,
		cartRepo.NewShoppingCartItemRepository,
		customerRepo.NewCustomerRepository,
		productRepo.NewProductRepository,
	)

	return &route.TransactionRoutes{}
}
