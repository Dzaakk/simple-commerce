//go:build wireinject
// +build wireinject

package injector

import (
	productRepo "Dzaakk/simple-commerce/internal/product/repositories"
	handler "Dzaakk/simple-commerce/internal/shopping_cart/handlers"
	repository "Dzaakk/simple-commerce/internal/shopping_cart/repositories"
	"Dzaakk/simple-commerce/internal/shopping_cart/routes"
	usecase "Dzaakk/simple-commerce/internal/shopping_cart/usecases"
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
