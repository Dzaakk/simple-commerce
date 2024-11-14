//go:build wireinject
// +build wireinject

package injector

import (
	customer "Dzaakk/simple-commerce/internal/customer/repositories"
	product "Dzaakk/simple-commerce/internal/product/repositories"
	cart "Dzaakk/simple-commerce/internal/shopping_cart/repositories"
	handlers "Dzaakk/simple-commerce/internal/transaction/handlers"
	repositories "Dzaakk/simple-commerce/internal/transaction/repositories"
	routes "Dzaakk/simple-commerce/internal/transaction/routes"
	usecases "Dzaakk/simple-commerce/internal/transaction/usecases"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *routes.TransactionRoutes {
	wire.Build(
		repositories.NewTransactionRepository,
		usecases.NewTransactionUseCase,
		handlers.NewTransactionHandler,
		routes.NewTransactionRoutes,
		cart.NewShoppingCartRepository,
		cart.NewShoppingCartItemRepository,
		customer.NewCustomerRepository,
		product.NewProductRepository,
	)

	return &routes.TransactionRoutes{}
}
