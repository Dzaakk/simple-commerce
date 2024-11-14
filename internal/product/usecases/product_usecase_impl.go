package usecase

import (
	model "Dzaakk/simple-commerce/internal/product/models"
	repo "Dzaakk/simple-commerce/internal/product/repositories"
	"fmt"
)

type ProductUseCaseImpl struct {
	repo repo.ProductRepository
}

func NewProductUseCase(repo repo.ProductRepository) ProductUseCase {
	return &ProductUseCaseImpl{repo}
}

func (p *ProductUseCaseImpl) Create(dataReq model.TProduct) (*model.ProductRes, error) {
	data, err := p.repo.Create(dataReq)
	if err != nil {
		return nil, err
	}

	newProduct := &model.ProductRes{
		Id:          fmt.Sprintf("%d", data.Id),
		ProductName: data.ProductName,
		Price:       fmt.Sprintf("%0.f", data.Price),
		Stock:       fmt.Sprintf("%d", data.Stock),
		CategoryId:  fmt.Sprintf("%d", data.CategoryId),
	}
	return newProduct, nil
}

func (p *ProductUseCaseImpl) FilterByPrice(price int) ([]*model.ProductRes, error) {
	panic("unimplemented")
}

func (p *ProductUseCaseImpl) Update(dataReq model.TProduct) error {
	err := p.repo.Update(dataReq)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductUseCaseImpl) FindByCategoryId(categoryId int) ([]*model.ProductRes, error) {
	listData, err := p.repo.FindByCategoryId(categoryId)
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

func (p *ProductUseCaseImpl) FindByName(productName string) (*model.ProductRes, error) {
	data, err := p.repo.FindByName(productName)
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
