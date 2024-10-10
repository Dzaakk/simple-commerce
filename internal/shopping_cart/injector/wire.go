//go:build wireinject
// +build wireinject

package injector

import (
	productRepo "Dzaakk/simple-commerce/internal/product/repository"
	handler "Dzaakk/simple-commerce/internal/shopping_cart/handler"
	repository "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"Dzaakk/simple-commerce/internal/shopping_cart/routes"
	usecase "Dzaakk/simple-commerce/internal/shopping_cart/usecase"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *routes.ShoppingCartRoutes {
	wire.Build(
		repository.NewShoppingCartRepository,
		repository.NewShoppingCartItemRepository,
		productRepo.NewProductRepository,
		usecase.NewShoppingCartUseCase,
		handler.NewShoppingCartHandler,
		routes.NewShoppingCartRoutes,
	)

	return &routes.ShoppingCartRoutes{}
}
