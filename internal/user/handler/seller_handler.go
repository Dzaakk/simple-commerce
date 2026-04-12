package handler

import (
	"Dzaakk/simple-commerce/package/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) FindSellerByName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.Error(response.NewAppError(http.StatusBadRequest, "invalid request data"))
		return
	}

	data, err := h.SellerService.FindByShopName(ctx, name)
	if err != nil {
		ctx.Error(err)
		return
	}

	if len(data) == 0 {
		ctx.Error(response.NewAppError(http.StatusNotFound, "seller not found"))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}
