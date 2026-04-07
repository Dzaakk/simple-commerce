package route

import (
	"Dzaakk/simple-commerce/internal/cart/handler"
	"Dzaakk/simple-commerce/internal/cart/repository"
	"Dzaakk/simple-commerce/internal/cart/service"
	catalogRepo "Dzaakk/simple-commerce/internal/catalog/repository"
	catalogService "Dzaakk/simple-commerce/internal/catalog/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type CartRoutes struct {
	Handler *handler.CartHandler
}

func NewCartRoutes(handler *handler.CartHandler) *CartRoutes {
	return &CartRoutes{Handler: handler}
}

func (cr *CartRoutes) Route(r *gin.RouterGroup) {
	cart := r.Group("/api/v1/cart")
	{
		cart.GET("", cr.Handler.GetCart)
		cart.POST("/items", cr.Handler.AddItem)
		cart.PUT("/items", cr.Handler.UpdateItem)
		cart.DELETE("/items/:product_id", cr.Handler.DeleteItem)
		cart.DELETE("/items", cr.Handler.ClearItems)
	}
}

func InitializedService(db *sql.DB) *CartRoutes {
	cartRepo := repository.NewCartRepository(db)
	cartItemRepo := repository.NewCartItemRepository(db)

	productRepo := catalogRepo.NewProductRepository(db)
	inventoryRepo := catalogRepo.NewInventoryRepository(db)

	productService := catalogService.NewProductService(productRepo)
	inventoryService := catalogService.NewInventoryService(inventoryRepo)

	cartService := service.NewCartService(cartRepo, cartItemRepo, productService, inventoryService)
	cartHandler := handler.NewCartHandler(cartService)

	return NewCartRoutes(cartHandler)
}
