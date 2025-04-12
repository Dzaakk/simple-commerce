//go:build wireinject
// +build wireinject

package injector

import (
	"Dzaakk/simple-commerce/internal/auth/handlers"
	"Dzaakk/simple-commerce/internal/auth/repositories"
	"Dzaakk/simple-commerce/internal/auth/routes"
	"Dzaakk/simple-commerce/internal/auth/usecases"
	customerRepo "Dzaakk/simple-commerce/internal/customer/repositories"
	customerUsecase "Dzaakk/simple-commerce/internal/customer/usecases"
	sellerRepo "Dzaakk/simple-commerce/internal/seller/repositories"
	sellerUsecase "Dzaakk/simple-commerce/internal/seller/usecases"
	shoppingCartRepo "Dzaakk/simple-commerce/internal/shopping_cart/repositories"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *routes.AuthRoutes {
	wire.Build(
		repositories.NewAuthRepository,
		sellerRepo.NewSellerRepository,
		customerRepo.NewCustomerRepository,
		shoppingCartRepo.NewShoppingCartRepository,
		usecases.NewAuthUseCase,
		customerUsecase.NewCustomerUseCase,
		sellerUsecase.NewSellerUseCase,
		handlers.NewAtuhHandler,
		routes.NewAuthRoutes,
	)

	return &routes.AuthRoutes{}
}
