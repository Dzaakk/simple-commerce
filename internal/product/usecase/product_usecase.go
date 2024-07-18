package product

import (
	model "Dzaakk/synapsis/internal/product/models"
	repo "Dzaakk/synapsis/internal/product/repository"
	"fmt"
)

type ProductUseCase interface {
	FindByCategoryId(categoryId int) ([]*model.ProductRes, error)
}
type ProductUseCaseImpl struct {
	repo repo.ProductRepository
}

func NewProductUseCase(repo repo.ProductRepository) ProductUseCase {
	return &ProductUseCaseImpl{repo}
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
