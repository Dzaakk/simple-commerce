//go:build wireinject
// +build wireinject

package injector

import (
	authRepo "Dzaakk/simple-commerce/internal/auth/repository"
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"
	"Dzaakk/simple-commerce/internal/product/handler"
	"Dzaakk/simple-commerce/internal/product/repository"
	"Dzaakk/simple-commerce/internal/product/route"
	"Dzaakk/simple-commerce/internal/product/usecase"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

func InitializedService(db *sql.DB, redis *redis.Client) *route.ProductRoutes {
	wire.Build(
		repository.NewProductRepository,
		authRepo.NewAuthCacheSellerRepository,
		usecase.NewProductUseCase,
		handler.NewProductHandler,
		route.NewProductRoutes,
		middleware.NewJWTSellerMiddleware,
	)

	return &route.ProductRoutes{}
}
