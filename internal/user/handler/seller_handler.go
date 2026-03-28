package handler

import (
	"Dzaakk/simple-commerce/package/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) FindSellerByName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, response.InvalidRequestData())
		return
	}

	data, err := h.SellerService.FindByShopName(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.InternalServerError(err.Error()))
		return
	}

	if len(data) == 0 {
		ctx.JSON(http.StatusNotFound, response.NotFound("seller not found"))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(data))
}
