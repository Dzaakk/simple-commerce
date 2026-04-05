package route

import (
	"Dzaakk/simple-commerce/internal/catalog/handler"
	"Dzaakk/simple-commerce/internal/catalog/repository"
	"Dzaakk/simple-commerce/internal/catalog/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type CatalogRoutes struct {
	Handler *handler.CatalogHandler
}

func NewCatalogRoutes(handler *handler.CatalogHandler) *CatalogRoutes {
	return &CatalogRoutes{
		Handler: handler,
	}
}

func (cr *CatalogRoutes) Route(r *gin.RouterGroup) {
	api := r.Group("/api/v1")

	product := api.Group("/product")
	{
		product.POST("", cr.Handler.CreateProduct)
		product.PUT("/:id", cr.Handler.UpdateProduct)
		product.DELETE("/:id", cr.Handler.DeleteProduct)
		product.GET("", cr.Handler.FindAllProducts)
		product.GET("/:id", cr.Handler.FindProductByID)
		product.PATCH("/:id/stock", cr.Handler.UpdateProductStock)
	}

	category := api.Group("/category")
	{
		category.POST("", cr.Handler.CreateCategory)
		category.GET("", cr.Handler.FindAllCategories)
		category.GET("/:id", cr.Handler.FindCategoryByID)
	}
}

func InitializedService(db *sql.DB) *CatalogRoutes {
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	productService := service.NewProductService(productRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	catalogHandler := handler.NewCatalogHandler(productService, categoryService)

	return NewCatalogRoutes(catalogHandler)
}
