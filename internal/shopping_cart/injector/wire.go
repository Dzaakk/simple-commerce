//go:build wireinject
// +build wireinject

package shopping_cart

import (
	productRepo "Dzaakk/simple-commerce/internal/product/repository"
	handler "Dzaakk/simple-commerce/internal/shopping_cart/handler"
	repository "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	usecase "Dzaakk/simple-commerce/internal/shopping_cart/usecase"
	"database/sql"

	"github.com/google/wire"
)

func InitializedService(db *sql.DB) *handler.ShoppingCarthandler {
	wire.Build(
		repository.NewShoppingCartRepository,
		repository.NewShoppingCartItemRepository,
		productRepo.NewProductRepository,
		usecase.NewShoppingCartUseCase,
		handler.NewShoppingCartHandler,
	)

	return &handler.ShoppingCarthandler{}
}
