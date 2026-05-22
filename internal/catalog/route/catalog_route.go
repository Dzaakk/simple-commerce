package route

import (
	"Dzaakk/simple-commerce/internal/catalog/handler"
	"Dzaakk/simple-commerce/internal/catalog/repository"
	"Dzaakk/simple-commerce/internal/catalog/service"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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

	apiV2 := r.Group("/api/v2")

	productV2 := apiV2.Group("/product")
	{
		productV2.GET("", cr.Handler.FindAllProductsV2)
		productV2.GET("/:id", cr.Handler.FindProductByIDV2)
	}
}

func InitializedService(db *sql.DB, redis *redis.Client) *CatalogRoutes {
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	productService := service.NewProductService(productRepo, redis)
	categoryService := service.NewCategoryService(categoryRepo)

	catalogHandler := handler.NewCatalogHandler(productService, categoryService)

	return NewCatalogRoutes(catalogHandler)
}
