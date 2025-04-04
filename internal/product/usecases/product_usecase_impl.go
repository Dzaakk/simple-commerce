package usecases

import (
	model "Dzaakk/simple-commerce/internal/product/models"
	repo "Dzaakk/simple-commerce/internal/product/repositories"
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
	sellerId, _ := strconv.ParseInt(dataReq.SellerId, 0, 64)
	categoryId, _ := strconv.ParseInt(dataReq.CategoryId, 0, 0)
	stock, _ := strconv.ParseInt(dataReq.Stock, 0, 0)
	newProduct := model.TProduct{
		ProductName: dataReq.ProductName,
		Price:       float32(price),
		Stock:       int(stock),
		CategoryId:  int(categoryId),
		SellerId:    int(sellerId),
	}
	data, err := p.repo.Create(ctx, newProduct)
	if err != nil {
		return nil, err
	}

	productRes := &model.ProductRes{
		Id:          fmt.Sprintf("%d", data.Id),
		ProductName: data.ProductName,
		Price:       fmt.Sprintf("%0.f", data.Price),
		Stock:       fmt.Sprintf("%d", data.Stock),
		CategoryId:  fmt.Sprintf("%d", data.CategoryId),
		SellerId:    fmt.Sprintf("%d", data.SellerId),
	}
	return productRes, nil
}

func (p *ProductUseCaseImpl) FilterByPrice(ctx context.Context, price int) ([]*model.ProductRes, error) {
	panic("unimplemented")
}

func (p *ProductUseCaseImpl) Update(ctx context.Context, dataReq model.ProductReq) error {
	price, _ := strconv.ParseFloat(dataReq.Price, 32)
	sellerId, _ := strconv.ParseInt(dataReq.SellerId, 0, 64)
	id, _ := strconv.ParseInt(dataReq.Id, 0, 64)
	categoryId, _ := strconv.ParseInt(dataReq.CategoryId, 0, 0)
	stock, _ := strconv.ParseInt(dataReq.Stock, 0, 0)
	updatedProduct := model.TProduct{
		Id:          int(id),
		ProductName: dataReq.ProductName,
		Price:       float32(price),
		Stock:       int(stock),
		CategoryId:  int(categoryId),
		SellerId:    int(sellerId),
	}

	_, err := p.repo.Update(ctx, updatedProduct)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductUseCaseImpl) FindByCategoryId(ctx context.Context, categoryId int) ([]*model.ProductRes, error) {
	listData, err := p.repo.FindProductByFilters(ctx, nil, &categoryId)
	if err != nil {
		return nil, err
	}
	var listProduct []*model.ProductRes
	for _, p := range listData {
		product := model.ProductRes{
			ProductName: p.ProductName,
			Price:       fmt.Sprintf("%0.f", p.Price),
			Stock:       fmt.Sprintf("%d", p.Stock),
			CategoryId:  fmt.Sprintf("%d", p.CategoryId),
		}
		listProduct = append(listProduct, &product)
	}

	return listProduct, nil
}

func (p *ProductUseCaseImpl) FindByName(ctx context.Context, productName string) (*model.ProductRes, error) {
	data, err := p.repo.FindByName(ctx, productName)
	if err != nil {
		return nil, err
	}
	return &model.ProductRes{
		Id:          fmt.Sprintf("%d", data.Id),
		ProductName: productName,
		Price:       fmt.Sprintf("%.0f", data.Price),
		Stock:       fmt.Sprintf("%d", data.Stock),
		CategoryId:  fmt.Sprintf("%d", data.CategoryId),
	}, nil
}
