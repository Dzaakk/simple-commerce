package api

import (
	"net/http"

	"Dzaakk/simple-commerce/internal/api/generated"
	"Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/internal/catalog/service"
	"Dzaakk/simple-commerce/package/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CatalogV2Server struct {
	ProductService service.ProductService
}

func NewCatalogV2Server(productService service.ProductService) *CatalogV2Server {
	return &CatalogV2Server{ProductService: productService}
}

func RegisterCatalogV2Routes(router gin.IRouter, productService service.ProductService) {
	generated.RegisterHandlersWithOptions(router, NewCatalogV2Server(productService), generated.GinServerOptions{
		ErrorHandler: func(ctx *gin.Context, err error, statusCode int) {
			ctx.Error(response.NewAppError(statusCode, "invalid request data"))
		},
	})
}

func (s *CatalogV2Server) ListProductsV2(ctx *gin.Context, params generated.ListProductsV2Params) {
	req, err := productQueryReqFromGenerated(params)
	if err != nil {
		ctx.Error(err)
		return
	}

	data, err := s.ProductService.FindAllCached(ctx.Request.Context(), req)
	if err != nil {
		ctx.Error(err)
		return
	}

	res, err := toGeneratedProductListResponse(data)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (s *CatalogV2Server) GetProductByIdV2(ctx *gin.Context, id generated.ProductIdPath) {
	data, err := s.ProductService.FindByIDCached(ctx.Request.Context(), id.String())
	if err != nil {
		ctx.Error(err)
		return
	}

	res, err := toGeneratedProductResponse(data)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func productQueryReqFromGenerated(params generated.ListProductsV2Params) (dto.ProductQueryReq, error) {
	var req dto.ProductQueryReq

	if params.CategoryId != nil {
		req.CategoryID = params.CategoryId
	}
	if params.SellerId != nil {
		sellerID := params.SellerId.String()
		req.SellerID = &sellerID
	}
	if params.MinPrice != nil {
		if *params.MinPrice < 0 {
			return req, response.NewAppError(http.StatusBadRequest, "invalid request data")
		}
		req.MinPrice = params.MinPrice
	}
	if params.MaxPrice != nil {
		if *params.MaxPrice < 0 {
			return req, response.NewAppError(http.StatusBadRequest, "invalid request data")
		}
		req.MaxPrice = params.MaxPrice
	}
	if params.Name != nil {
		req.Name = params.Name
	}
	if params.Cursor != nil {
		req.Cursor = params.Cursor
	}
	if params.Limit != nil {
		if *params.Limit < 1 {
			return req, response.NewAppError(http.StatusBadRequest, "invalid request data")
		}
		req.Limit = *params.Limit
	}
	if params.SortBy != nil {
		switch *params.SortBy {
		case generated.ListProductsV2ParamsSortByNewest,
			generated.ListProductsV2ParamsSortByPriceAsc,
			generated.ListProductsV2ParamsSortByPriceDesc:
			req.SortBy = string(*params.SortBy)
		default:
			return req, response.NewAppError(http.StatusBadRequest, "invalid request data")
		}
	}

	return req, nil
}

func toGeneratedProductListResponse(data *dto.ProductListRes) (generated.ProductListResponse, error) {
	items := make([]generated.Product, 0)
	var nextCursor *string
	if data != nil {
		items = make([]generated.Product, 0, len(data.Items))
		for _, item := range data.Items {
			product, err := toGeneratedProduct(item)
			if err != nil {
				return generated.ProductListResponse{}, err
			}
			items = append(items, product)
		}
		nextCursor = data.NextCursor
	}

	return generated.ProductListResponse{
		Meta: successMeta(),
		Data: generated.ProductList{
			Items:      items,
			NextCursor: nextCursor,
		},
	}, nil
}

func toGeneratedProductResponse(data *dto.ProductRes) (generated.ProductResponse, error) {
	if data == nil {
		return generated.ProductResponse{}, response.NewAppError(http.StatusNotFound, "product not found")
	}

	product, err := toGeneratedProduct(*data)
	if err != nil {
		return generated.ProductResponse{}, err
	}

	return generated.ProductResponse{
		Meta: successMeta(),
		Data: product,
	}, nil
}

func toGeneratedProduct(data dto.ProductRes) (generated.Product, error) {
	productID, err := uuid.Parse(data.ID)
	if err != nil {
		return generated.Product{}, response.NewAppError(http.StatusInternalServerError, "internal server error")
	}

	var sellerID *uuid.UUID
	if data.SellerID != "" {
		parsedSellerID, err := uuid.Parse(data.SellerID)
		if err != nil {
			return generated.Product{}, response.NewAppError(http.StatusInternalServerError, "internal server error")
		}
		sellerID = &parsedSellerID
	}

	return generated.Product{
		Id:          productID,
		SellerId:    sellerID,
		CategoryId:  &data.CategoryID,
		Name:        &data.Name,
		Sku:         &data.SKU,
		Description: data.Description,
		Price:       &data.Price,
		ImageUrl:    data.ImageURL,
		IsActive:    &data.IsActive,
		CreatedAt:   &data.CreatedAt,
		UpdatedAt:   &data.UpdatedAt,
	}, nil
}

func successMeta() generated.Meta {
	return generated.Meta{
		Code:    http.StatusOK,
		Message: "Success",
	}
}
