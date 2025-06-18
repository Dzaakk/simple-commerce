package route

import (
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"
	handler "Dzaakk/simple-commerce/internal/shopping_cart/handler"

	"github.com/gin-gonic/gin"
)

type ShoppingCartRoutes struct {
	Handler       *handler.ShoppingCartHandler
	JWTMiddleware *middleware.JWTCustomerMiddleware
}

func NewShoppingCartRoutes(handler *handler.ShoppingCartHandler, jwtMiddleware *middleware.JWTCustomerMiddleware) *ShoppingCartRoutes {
	return &ShoppingCartRoutes{
		Handler:       handler,
		JWTMiddleware: jwtMiddleware,
	}
}
func (scr *ShoppingCartRoutes) Route(r *gin.RouterGroup) {
	ShoppingHandler := r.Group("api/v1")

	ShoppingHandler.Use(scr.JWTMiddleware.ValidateToken())
	{
		ShoppingHandler.POST("/shopping-cart", scr.Handler.AddProductToShoppingCart)
		ShoppingHandler.GET("/shopping-cart", scr.Handler.GetListShoppingCart)
		ShoppingHandler.POST("/shopping-cart/delete", scr.Handler.DeleteShoppingList)
	}
}
