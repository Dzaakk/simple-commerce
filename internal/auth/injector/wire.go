//go:build wireinject
// +build wireinject

package injector

import (
	"Dzaakk/simple-commerce/internal/auth/handler"
	"Dzaakk/simple-commerce/internal/auth/repository"
	"Dzaakk/simple-commerce/internal/auth/route"
	"Dzaakk/simple-commerce/internal/auth/usecase"
	customerRepo "Dzaakk/simple-commerce/internal/customer/repository"
	sellerRepo "Dzaakk/simple-commerce/internal/seller/repository"
	sellerUsecase "Dzaakk/simple-commerce/internal/seller/usecase"
	shoppingCartRepo "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"database/sql"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

func InitializedService(db *sql.DB, redis *redis.Client) *route.AuthRoutes {
	wire.Build(
		repository.NewAuthCacheRepository,
		sellerRepo.NewSellerRepository,
		customerRepo.NewCustomerRepository,
		shoppingCartRepo.NewShoppingCartRepository,
		usecase.NewAuthUseCase,
		sellerUsecase.NewSellerUseCase,
		handler.NewAtuhHandler,
		route.NewAuthRoutes,
	)

	return &route.AuthRoutes{}
}
