package route

import (
	middleware "Dzaakk/simple-commerce/internal/middleware/jwt"
	"Dzaakk/simple-commerce/internal/product/handler"

	"github.com/gin-gonic/gin"
)

type ProductRoutes struct {
	Handler          *handler.ProductHandler
	SellerMiddleware *middleware.JWTSellerMiddleware
}

func NewProductRoutes(handler *handler.ProductHandler, sellerMiddleware *middleware.JWTSellerMiddleware) *ProductRoutes {
	return &ProductRoutes{
		Handler:          handler,
		SellerMiddleware: sellerMiddleware,
	}
}
func (pr *ProductRoutes) Route(r *gin.RouterGroup) {
	productHandler := r.Group("api/v1")

	productHandler.Use()
	{
		productHandler.POST("/product", pr.SellerMiddleware.ValidateToken(), pr.Handler.CreateProduct)
		productHandler.GET("/product", pr.Handler.GetProducts)
		productHandler.POST("/product", pr.SellerMiddleware.ValidateToken(), pr.Handler.UpdateProduct)
		productHandler.POST("/product", pr.SellerMiddleware.ValidateToken(), pr.Handler.DeleteProduct)
	}
}
