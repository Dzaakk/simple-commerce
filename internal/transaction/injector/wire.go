//go:build wireinject
// +build wireinject

package injector

import (
	auth "Dzaakk/simple-commerce/internal/auth/repository"
	customerRepo "Dzaakk/simple-commerce/internal/customer/repository"
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"
	productRepo "Dzaakk/simple-commerce/internal/product/repository"
	cartRepo "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"Dzaakk/simple-commerce/internal/transaction/handler"
	"Dzaakk/simple-commerce/internal/transaction/repository"
	"Dzaakk/simple-commerce/internal/transaction/route"
	"Dzaakk/simple-commerce/internal/transaction/usecase"
	"database/sql"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

func InitializedService(db *sql.DB, redis *redis.Client) *route.TransactionRoutes {
	wire.Build(
		repository.NewTransactionRepository,
		usecase.NewTransactionUseCase,
		handler.NewTransactionHandler,
		route.NewTransactionRoutes,
		cartRepo.NewShoppingCartRepository,
		cartRepo.NewShoppingCartItemRepository,
		customerRepo.NewCustomerRepository,
		productRepo.NewProductRepository,
		auth.NewAuthCacheRepository,
		middleware.NewJwtMiddleware,
	)

	return &route.TransactionRoutes{}
}
