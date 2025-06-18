//go:build wireinject
// +build wireinject

package injector

import (
	auth "Dzaakk/simple-commerce/internal/auth/repository"
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"
	productRepo "Dzaakk/simple-commerce/internal/product/repository"
	"Dzaakk/simple-commerce/internal/shopping_cart/handler"
	"Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"Dzaakk/simple-commerce/internal/shopping_cart/route"
	"Dzaakk/simple-commerce/internal/shopping_cart/usecase"
	"database/sql"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

func InitializedService(db *sql.DB, redis *redis.Client) *route.ShoppingCartRoutes {
	wire.Build(
		repository.NewShoppingCartRepository,
		repository.NewShoppingCartItemRepository,
		productRepo.NewProductRepository,
		usecase.NewShoppingCartUseCase,
		handler.NewShoppingCartHandler,
		route.NewShoppingCartRoutes,
		auth.NewAuthCacheCustomerRepository,
		middleware.NewJWTCustomerMiddleware,
	)

	return &route.ShoppingCartRoutes{}
}
