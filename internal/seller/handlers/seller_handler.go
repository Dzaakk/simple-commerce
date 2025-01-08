package handlers

import (
	usecase "Dzaakk/simple-commerce/internal/seller/usecases"

	"github.com/gin-gonic/gin"
)

type SellerHandler struct {
	Usecase usecase.SellerUseCase
}

func NewSellerHandler(usecase usecase.SellerUseCase) *SellerHandler {
	return &SellerHandler{Usecase: usecase}
}

func (handler *SellerHandler) Register(ctx *gin.Context) {

}
func (handler *SellerHandler) Login(ctx *gin.Context) {

}
func (handler *SellerHandler) FindById(ctx *gin.Context) {

}
func (handler *SellerHandler) Deactivate(ctx *gin.Context) {

}
func (handler *SellerHandler) ChangePassword(ctx *gin.Context) {

}
