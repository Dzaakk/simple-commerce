//go:build wireinject
// +build wireinject

package injector

import (
	auth "Dzaakk/simple-commerce/internal/auth/repository"
	"Dzaakk/simple-commerce/internal/customer/handler"
	"Dzaakk/simple-commerce/internal/customer/repository"
	"Dzaakk/simple-commerce/internal/customer/route"
	"Dzaakk/simple-commerce/internal/customer/usecase"
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"
	"database/sql"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

func InitializedService(db *sql.DB, redis *redis.Client) *route.CustomerRoutes {
	wire.Build(
		repository.NewCustomerRepository,
		usecase.NewCustomerUseCase,
		handler.NewCustomerHandler,
		auth.NewAuthCacheRepository,
		middleware.NewJwtMiddleware,
		route.NewCustomerRoutes,
	)

	return &route.CustomerRoutes{}
}
