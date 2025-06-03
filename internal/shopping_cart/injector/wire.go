//go:build wireinject
// +build wireinject

package injector

import (
	productRepo "Dzaakk/simple-commerce/internal/product/repository"
	"Dzaakk/simple-commerce/internal/shopping_cart/handler"
	"Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"Dzaakk/simple-commerce/internal/shopping_cart/route"
	"Dzaakk/simple-commerce/internal/shopping_cart/usecase"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *route.ShoppingCartRoutes {
	wire.Build(
		repository.NewShoppingCartRepository,
		repository.NewShoppingCartItemRepository,
		productRepo.NewProductRepository,
		usecase.NewShoppingCartUseCase,
		handler.NewShoppingCartHandler,
		route.NewShoppingCartRoutes,
	)

	return &route.ShoppingCartRoutes{}
}
