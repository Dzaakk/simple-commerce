package handler

import (
	"Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/package/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func productQueryReqFromContext(ctx *gin.Context) (dto.ProductQueryReq, error) {
	var req dto.ProductQueryReq

	if val := ctx.Query("category_id"); val != "" {
		id, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return req, response.NewAppError(http.StatusBadRequest, "invalid request data")
		}
		req.CategoryID = &id
	}
	if val := ctx.Query("seller_id"); val != "" {
		req.SellerID = &val
	}
	if val := ctx.Query("min_price"); val != "" {
		minPrice, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return req, response.NewAppError(http.StatusBadRequest, "invalid request data")
		}
		req.MinPrice = &minPrice
	}
	if val := ctx.Query("max_price"); val != "" {
		maxPrice, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return req, response.NewAppError(http.StatusBadRequest, "invalid request data")
		}
		req.MaxPrice = &maxPrice
	}
	if val := ctx.Query("name"); val != "" {
		req.Name = &val
	}
	if val := ctx.Query("cursor"); val != "" {
		req.Cursor = &val
	}
	if val := ctx.Query("limit"); val != "" {
		limit, err := strconv.Atoi(val)
		if err != nil {
			return req, response.NewAppError(http.StatusBadRequest, "invalid request data")
		}
		req.Limit = limit
	}
	req.SortBy = ctx.Query("sort_by")

	return req, nil
}
