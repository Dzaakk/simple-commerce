package handler

import (
	usecase "Dzaakk/simple-commerce/internal/product/usecase"
	template "Dzaakk/simple-commerce/package/template"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type ProductHandler struct {
	Usecase usecase.ProductUseCase
}

func NewProductHandler(usecase usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		Usecase: usecase,
	}
}

func (handler *ProductHandler) Route(r *gin.RouterGroup, redis *redis.Client) {
	productHandler := r.Group("api/v1")

	productHandler.Use()
	{
		productHandler.GET("/product", func(ctx *gin.Context) {
			if err := handler.GetProduct(ctx); err != nil {
				ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "failed to Get product", err.Error()))
				return
			}
		})
	}
}

func (handler *ProductHandler) GetProduct(ctx *gin.Context) error {
	id, _ := strconv.Atoi(ctx.Query("id"))
	categoryId, _ := strconv.Atoi(ctx.Query("categoryId"))

	if id != 0 {

	}

	if categoryId != 0 {
		listProduct, err := handler.Usecase.FindByCategoryId(categoryId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, template.Response(http.StatusInternalServerError, "Internal Server Error", err.Error()))
			return nil
		}
		if listProduct == nil {
			ctx.JSON(http.StatusNotFound, template.Response(http.StatusNotFound, "not found", fmt.Sprintf("product with category %v is not found", categoryId)))
			ctx.Abort()
			return nil
		}
		ctx.JSON(http.StatusOK, template.Response(http.StatusOK, "Success", listProduct))
		return nil
	}

	return nil
}
