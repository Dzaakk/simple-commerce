package route

import (
	handler "Dzaakk/simple-commerce/internal/shopping_cart/handler"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type ShoppingCartRoutes struct {
	Handler *handler.ShoppingCartHandler
}

func NewShoppingCartRoutes(handler *handler.ShoppingCartHandler) *ShoppingCartRoutes {
	return &ShoppingCartRoutes{
		Handler: handler,
	}
}
func (scr *ShoppingCartRoutes) Route(r *gin.RouterGroup, redis *redis.Client) {
	ShoppingHandler := r.Group("api/v1")

	ShoppingHandler.Use()
	{
		ShoppingHandler.POST("/shopping-cart", scr.Handler.AddProductToShoppingCart)
		ShoppingHandler.GET("/shopping-cart", scr.Handler.GetListShoppingCart)
		ShoppingHandler.POST("/shopping-cart/delete", scr.Handler.DeleteShoppingList)
	}
}
