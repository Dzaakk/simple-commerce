package usecase

import (
	"Dzaakk/simple-commerce/internal/product/model"
	repo "Dzaakk/simple-commerce/internal/product/repository"
	"context"
	"fmt"
	"strconv"
)

type ProductUseCaseImpl struct {
	repo repo.ProductRepository
}

func NewProductUseCase(repo repo.ProductRepository) ProductUseCase {
	return &ProductUseCaseImpl{repo}
}

func (p *ProductUseCaseImpl) Create(ctx context.Context, dataReq model.ProductReq) (*model.ProductRes, error) {
	price, _ := strconv.ParseFloat(dataReq.Price, 32)
	sellerID, _ := strconv.ParseInt(dataReq.SellerID, 0, 64)
	categoryID, _ := strconv.ParseInt(dataReq.CategoryID, 0, 0)
	stock, _ := strconv.ParseInt(dataReq.Stock, 0, 0)
	newProduct := model.TProduct{
		ProductName: dataReq.ProductName,
		Price:       float32(price),
		Stock:       int(stock),
		CategoryID:  int(categoryID),
		SellerID:    int(sellerID),
	}
	data, err := p.repo.Create(ctx, newProduct)
	if err != nil {
		return nil, err
	}

	productRes := &model.ProductRes{
		ProductID:   fmt.Sprintf("%d", data.ID),
		ProductName: data.ProductName,
		Price:       fmt.Sprintf("%0.f", data.Price),
		Stock:       fmt.Sprintf("%d", data.Stock),
		CategoryID:  fmt.Sprintf("%d", data.CategoryID),
		SellerID:    fmt.Sprintf("%d", data.SellerID),
	}
	return productRes, nil
}

func (p *ProductUseCaseImpl) Update(ctx context.Context, dataReq model.ProductReq) error {
	price, _ := strconv.ParseFloat(dataReq.Price, 32)
	sellerID, _ := strconv.ParseInt(dataReq.SellerID, 0, 64)
	productID, _ := strconv.ParseInt(dataReq.ProductID, 0, 64)
	categoryID, _ := strconv.ParseInt(dataReq.CategoryID, 0, 0)
	stock, _ := strconv.ParseInt(dataReq.Stock, 0, 0)
	updatedProduct := model.TProduct{
		ID:          int(productID),
		ProductName: dataReq.ProductName,
		Price:       float32(price),
		Stock:       int(stock),
		CategoryID:  int(categoryID),
		SellerID:    int(sellerID),
	}

	_, err := p.repo.Update(ctx, updatedProduct)
	if err != nil {
		return err
	}

	return nil
}
func (p *ProductUseCaseImpl) FindByFilter(ctx context.Context, params model.ProductFilter) ([]*model.ProductRes, error) {

	listProduct, err := p.repo.FindByFilters(ctx, params)
	if err != nil {
		return nil, err
	}

	var listData []*model.ProductRes
	for _, p := range listProduct {
		product := p.ToResponse()
		listData = append(listData, &product)
	}

	return listData, nil
}

func (p *ProductUseCaseImpl) FindByCategoryID(ctx context.Context, categoryID int) ([]*model.ProductRes, error) {
	listData, err := p.repo.FindProductByFilters(ctx, nil, &categoryID)
	if err != nil {
		return nil, err
	}
	var listProduct []*model.ProductRes
	for _, p := range listData {
		product := model.ProductRes{
			ProductName: p.ProductName,
			Price:       fmt.Sprintf("%0.f", p.Price),
			Stock:       fmt.Sprintf("%d", p.Stock),
			CategoryID:  fmt.Sprintf("%d", p.CategoryID),
		}
		listProduct = append(listProduct, &product)
	}

	return listProduct, nil
}

func (p *ProductUseCaseImpl) FindByProductName(ctx context.Context, productName string) (*model.ProductRes, error) {
	data, err := p.repo.FindByProductName(ctx, productName)
	if err != nil {
		return nil, err
	}
	return &model.ProductRes{
		ProductID:   fmt.Sprintf("%d", data.ID),
		ProductName: productName,
		Price:       fmt.Sprintf("%.0f", data.Price),
		Stock:       fmt.Sprintf("%d", data.Stock),
		CategoryID:  fmt.Sprintf("%d", data.CategoryID),
	}, nil
}
