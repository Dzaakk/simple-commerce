package handler

import (
	"Dzaakk/simple-commerce/internal/user/service"
	"Dzaakk/simple-commerce/package/response"
	"Dzaakk/simple-commerce/package/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	CustomerService service.CustomerService
	SellerService   service.SellerService
}

func NewUserHandler(customerService service.CustomerService, sellerService service.SellerService) *UserHandler {
	return &UserHandler{
		CustomerService: customerService,
		SellerService:   sellerService,
	}
}

func (h *UserHandler) FindCustomerByID(ctx *gin.Context) {

	customerID := ctx.Query("id")
	if !util.AuthorizedChecker(ctx, ctx.Query("id")) {
		ctx.Error(response.NewAppError(http.StatusUnauthorized, "unauthorized"))
		return
	}

	data, err := h.CustomerService.FindByID(ctx, customerID)
	if err != nil {
		ctx.Error(err)
		return
	}

	if data == nil {
		ctx.Error(response.NewAppError(http.StatusNotFound, "user not found"))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}

func (h *UserHandler) FindCustomerByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	data, err := h.CustomerService.FindByEmail(ctx, email)
	if err != nil {
		ctx.Error(err)
		return
	}

	if data == nil {
		ctx.Error(response.NewAppError(http.StatusNotFound, "user not found"))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}
