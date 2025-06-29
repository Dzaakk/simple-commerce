// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package injector

import (
	repository3 "Dzaakk/simple-commerce/internal/auth/repository"
	"Dzaakk/simple-commerce/internal/middleware/jwt"
	repository2 "Dzaakk/simple-commerce/internal/product/repository"
	"Dzaakk/simple-commerce/internal/shopping_cart/handler"
	"Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"Dzaakk/simple-commerce/internal/shopping_cart/route"
	"Dzaakk/simple-commerce/internal/shopping_cart/usecase"
	"database/sql"
	"github.com/go-redis/redis/v8"
)

// Injectors from wire.go:

func InitializedService(db *sql.DB, redis2 *redis.Client) *route.ShoppingCartRoutes {
	shoppingCartRepository := repository.NewShoppingCartRepository(db)
	shoppingCartItemRepository := repository.NewShoppingCartItemRepository(db)
	productRepository := repository2.NewProductRepository(db)
	shoppingCartUseCase := usecase.NewShoppingCartUseCase(shoppingCartRepository, shoppingCartItemRepository, productRepository)
	shoppingCartHandler := handler.NewShoppingCartHandler(shoppingCartUseCase)
	authCacheCustomer := repository3.NewAuthCacheCustomerRepository(redis2)
	jwtCustomerMiddleware := middleware.NewJWTCustomerMiddleware(authCacheCustomer)
	shoppingCartRoutes := route.NewShoppingCartRoutes(shoppingCartHandler, jwtCustomerMiddleware)
	return shoppingCartRoutes
}
