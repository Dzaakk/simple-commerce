//go:build wireinject
// +build wireinject

package injector

import (
	"Dzaakk/simple-commerce/internal/auth/handler"
	"Dzaakk/simple-commerce/internal/auth/repository"
	"Dzaakk/simple-commerce/internal/auth/route"
	"Dzaakk/simple-commerce/internal/auth/usecase"
	customerRepo "Dzaakk/simple-commerce/internal/customer/repository"
	emailUsecase "Dzaakk/simple-commerce/internal/email/usecase"
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"
	sellerRepo "Dzaakk/simple-commerce/internal/seller/repository"
	shoppingCartRepo "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

func InitializedService(db *sql.DB, redis *redis.Client) *route.AuthRoutes {
	wire.Build(
		repository.NewAuthCacheSellerRepository,
		repository.NewAuthCacheCustomerRepository,
		sellerRepo.NewSellerRepository,
		customerRepo.NewCustomerRepository,
		shoppingCartRepo.NewShoppingCartRepository,
		usecase.NewAuthUseCase,
		handler.NewAtuhHandler,
		middleware.NewJWTCustomerMiddleware,
		middleware.NewJWTSellerMiddleware,
		route.NewAuthRoutes,
		emailUsecase.NewEmailUseCase,
	)

	return &route.AuthRoutes{}
}
