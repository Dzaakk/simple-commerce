package service

import (
	"Dzaakk/simple-commerce/internal/catalog/dto"
	repo "Dzaakk/simple-commerce/internal/catalog/repository"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"net/http"
	"strconv"
	"time"
)

type ProductServiceImpl struct {
	Repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return &ProductServiceImpl{Repo: repo}
}

func (p *ProductServiceImpl) Create(ctx context.Context, req *dto.CreateProductReq) (string, error) {
	data := req.ToCreateData()

	id, err := p.Repo.Create(ctx, data)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (p *ProductServiceImpl) Update(ctx context.Context, productID string, sellerID string, req *dto.UpdateProductReq) error {
	if productID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter product id")
	}
	if sellerID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter seller id")
	}

	data := req.ToUpdateData(productID, sellerID)

	rowsAffected, err := p.Repo.Update(ctx, data)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return response.NewAppError(http.StatusNotFound, "product not found")
	}

	return nil
}

func (p *ProductServiceImpl) SoftDelete(ctx context.Context, productID string, sellerID string) error {
	if productID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter product id")
	}
	if sellerID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter seller id")
	}

	rowsAffected, err := p.Repo.SoftDelete(ctx, productID, time.Now())
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return response.NewAppError(http.StatusNotFound, "product not found")
	}

	return nil
}

func (p *ProductServiceImpl) FindByID(ctx context.Context, productID string) (*dto.ProductRes, error) {
	data, err := p.Repo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, response.NewAppError(http.StatusNotFound, "product not found")
	}

	product := dto.ToProductRes(data)

	return &product, nil
}

func (p *ProductServiceImpl) FindAll(ctx context.Context, req dto.ProductQueryReq) (*dto.ProductListRes, error) {
	filter := repo.ProductFilter{
		CategoryID: req.CategoryID,
		SellerID:   req.SellerID,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Name:       req.Name,
		Cursor:     req.Cursor,
		Limit:      req.Limit,
		SortBy:     req.SortBy,
	}

	data, err := p.Repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return &dto.ProductListRes{Items: []dto.ProductRes{}}, nil
	}

	result := make([]dto.ProductRes, 0, len(data))
	for _, product := range data {
		if product == nil {
			continue
		}
		res := dto.ToProductRes(product)
		result = append(result, res)
	}

	res := &dto.ProductListRes{Items: result}

	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "newest"
	}

	if req.Limit > 0 && len(result) == req.Limit {
		cursor := buildProductCursor(sortBy, result[len(result)-1])
		res.NextCursor = &cursor
	}

	return res, nil
}

func (p *ProductServiceImpl) UpdateStock(ctx context.Context, productID string, sellerID string, quantity int) error {
	if productID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter product id")
	}
	if sellerID == "" {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter seller id")
	}
	if quantity < 0 {
		return response.NewAppError(http.StatusBadRequest, "invalid parameter quantity")
	}

	return p.Repo.UpdateStock(ctx, productID, sellerID, quantity)
}

func buildProductCursor(sortBy string, p dto.ProductRes) string {
	switch sortBy {
	case "price_asc", "price_desc":
		return strconv.FormatFloat(p.Price, 'f', -1, 64) + "|" + p.ID
	default:
		return p.CreatedAt.Format(time.RFC3339Nano) + "|" + p.ID
	}
}
