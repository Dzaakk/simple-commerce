package handlers

import (
	usecase "Dzaakk/simple-commerce/internal/seller/usecases"
	"Dzaakk/simple-commerce/package/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SellerHandler struct {
	Usecase usecase.SellerUseCase
}

func NewSellerHandler(usecase usecase.SellerUseCase) *SellerHandler {
	return &SellerHandler{Usecase: usecase}
}

func (handler *SellerHandler) FindById(ctx *gin.Context) {

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

}
func (handler *SellerHandler) Deactivate(ctx *gin.Context) {

}
func (handler *SellerHandler) ChangePassword(ctx *gin.Context) {

}
