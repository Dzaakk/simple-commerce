package handler

import (
	"Dzaakk/simple-commerce/internal/seller/usecase"
	"Dzaakk/simple-commerce/package/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SellerHandler struct {
	Usecase usecase.SellerUseCase
}

func NewSellerHandler(usecase usecase.SellerUseCase) *SellerHandler {
	return &SellerHandler{Usecase: usecase}
}

func (handler *SellerHandler) FindById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("id"))

	if id != 0 {
		seller, err := handler.Usecase.FindById(ctx, int64(id))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, response.Success(seller))
	}
}

func (handler *SellerHandler) FindByUsername(ctx *gin.Context) {
	username := ctx.Query("username")
	if username != "" {
		seller, err := handler.Usecase.FindByUsername(ctx, username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, response.Success(seller))
	}
}
func (handler *SellerHandler) FindAll(ctx *gin.Context) {
	seller, err := handler.Usecase.FindAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.Success(seller))
}
