//go:build wireinject
// +build wireinject

package injector

import (
	"Dzaakk/simple-commerce/internal/auth/handler"
	"Dzaakk/simple-commerce/internal/auth/repository"
	"Dzaakk/simple-commerce/internal/auth/route"
	"Dzaakk/simple-commerce/internal/auth/usecase"
	customerRepo "Dzaakk/simple-commerce/internal/customer/repository"
	customerUsecase "Dzaakk/simple-commerce/internal/customer/usecase"
	sellerRepo "Dzaakk/simple-commerce/internal/seller/repository"
	sellerUsecase "Dzaakk/simple-commerce/internal/seller/usecase"
	shoppingCartRepo "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *route.AuthRoutes {
	wire.Build(
		repository.NewAuthRepository,
		sellerRepo.NewSellerRepository,
		customerRepo.NewCustomerRepository,
		shoppingCartRepo.NewShoppingCartRepository,
		usecase.NewAuthUseCase,
		customerUsecase.NewCustomerUseCase,
		sellerUsecase.NewSellerUseCase,
		handler.NewAtuhHandler,
		route.NewAuthRoutes,
	)

	return &route.AuthRoutes{}
}
