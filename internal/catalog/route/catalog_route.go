package route

import (
	contractapi "Dzaakk/simple-commerce/internal/api"
	"Dzaakk/simple-commerce/internal/catalog/handler"
	"Dzaakk/simple-commerce/internal/catalog/repository"
	"Dzaakk/simple-commerce/internal/catalog/service"
	"Dzaakk/simple-commerce/internal/middleware"
	"Dzaakk/simple-commerce/package/constant"
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
		product.GET("", cr.Handler.FindAllProducts)
		product.GET("/:id", cr.Handler.FindProductByID)

		sellerProduct := product.Group("", middleware.Authenticate(), middleware.RequireUserType(constant.Seller))
		sellerProduct.POST("", cr.Handler.CreateProduct)
		sellerProduct.PUT("/:id", cr.Handler.UpdateProduct)
		sellerProduct.DELETE("/:id", cr.Handler.DeleteProduct)
		sellerProduct.PATCH("/:id/stock", cr.Handler.UpdateProductStock)
	}

	category := api.Group("/category")
	{
		category.POST("", cr.Handler.CreateCategory)
		category.GET("", cr.Handler.FindAllCategories)
		category.GET("/:id", cr.Handler.FindCategoryByID)
	}

	contractapi.RegisterCatalogV2Routes(r, cr.Handler.ProductService)
}

func InitializedService(db *sql.DB, redis *redis.Client) *CatalogRoutes {
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	productService := service.NewProductService(productRepo, redis)
	categoryService := service.NewCategoryService(categoryRepo)

	catalogHandler := handler.NewCatalogHandler(productService, categoryService)

	return NewCatalogRoutes(catalogHandler)
}
