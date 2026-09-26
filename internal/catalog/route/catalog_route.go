package route

import (
	"database/sql"

	"Dzaakk/simple-commerce/internal/catalog/handler"
	"Dzaakk/simple-commerce/internal/catalog/repository"
	"Dzaakk/simple-commerce/internal/catalog/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func Route(router *gin.RouterGroup, db *sql.DB, redisClient *redis.Client) {
	products := service.NewProductService(repository.NewProductRepository(db), redisClient)
	categories := service.NewCategoryService(repository.NewCategoryRepository(db))
	h := handler.NewCatalogHandler(products, categories)

	router.GET("/api/v1/product", h.FindAllProducts)
	router.GET("/api/v1/product/:id", h.FindProductByID)
	router.GET("/api/v1/category", h.FindAllCategories)
	router.GET("/api/v2/product", h.FindAllProductsV2)
	router.GET("/api/v2/product/:id", h.FindProductByIDV2)
}
