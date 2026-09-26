package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/package/response"

	"github.com/gin-gonic/gin"
)

func productQuery(ctx *gin.Context) (dto.ProductQuery, error) {
	query := dto.ProductQuery{Limit: 20, SortBy: ctx.DefaultQuery("sort_by", "newest")}
	if query.SortBy != "newest" && query.SortBy != "price_asc" && query.SortBy != "price_desc" {
		return query, response.NewAppError(http.StatusBadRequest, "invalid sort_by")
	}
	if value := ctx.Query("category_id"); value != "" {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id <= 0 {
			return query, response.NewAppError(http.StatusBadRequest, "invalid category_id")
		}
		query.CategoryID = &id
	}
	if value := ctx.Query("min_price"); value != "" {
		price, err := strconv.ParseFloat(value, 64)
		if err != nil || price < 0 {
			return query, response.NewAppError(http.StatusBadRequest, "invalid min_price")
		}
		query.MinPrice = &price
	}
	if value := ctx.Query("max_price"); value != "" {
		price, err := strconv.ParseFloat(value, 64)
		if err != nil || price < 0 {
			return query, response.NewAppError(http.StatusBadRequest, "invalid max_price")
		}
		query.MaxPrice = &price
	}
	if query.MinPrice != nil && query.MaxPrice != nil && *query.MinPrice > *query.MaxPrice {
		return query, response.NewAppError(http.StatusBadRequest, "min_price cannot exceed max_price")
	}
	if value := strings.TrimSpace(ctx.Query("name")); value != "" {
		query.Name = &value
	}
	if value := ctx.Query("cursor"); value != "" {
		if !validCursor(value, query.SortBy) {
			return query, response.NewAppError(http.StatusBadRequest, "invalid cursor")
		}
		query.Cursor = &value
	}
	if value := ctx.Query("limit"); value != "" {
		limit, err := strconv.Atoi(value)
		if err != nil || limit < 1 || limit > 100 {
			return query, response.NewAppError(http.StatusBadRequest, "limit must be between 1 and 100")
		}
		query.Limit = limit
	}
	return query, nil
}

func validCursor(cursor, sortBy string) bool {
	parts := strings.SplitN(cursor, "|", 2)
	if len(parts) != 2 {
		return false
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || id <= 0 {
		return false
	}
	if sortBy == "price_asc" || sortBy == "price_desc" {
		price, err := strconv.ParseFloat(parts[0], 64)
		return err == nil && price >= 0
	}
	_, err = time.Parse(time.RFC3339Nano, parts[0])
	return err == nil
}
