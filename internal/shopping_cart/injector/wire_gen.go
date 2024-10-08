// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package injector

import (
	repository2 "Dzaakk/simple-commerce/internal/product/repository"
	"Dzaakk/simple-commerce/internal/shopping_cart/handler"
	"Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"Dzaakk/simple-commerce/internal/shopping_cart/routes"
	"Dzaakk/simple-commerce/internal/shopping_cart/usecase"
	"database/sql"
)

// Injectors from wire.go:

func InitializedService(db *sql.DB) *routes.ShoppingCartRoutes {
	shoppingCartRepository := repository.NewShoppingCartRepository(db)
	shoppingCartItemRepository := repository.NewShoppingCartItemRepository(db)
	productRepository := repository2.NewProductRepository(db)
	shoppingCartUseCase := usecase.NewShoppingCartUseCase(shoppingCartRepository, shoppingCartItemRepository, productRepository)
	shoppingCartHandler := handler.NewShoppingCartHandler(shoppingCartUseCase)
	shoppingCartRoutes := routes.NewShoppingCartRoutes(shoppingCartHandler)
	return shoppingCartRoutes
}
