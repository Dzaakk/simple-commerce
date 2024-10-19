package usecase

import (
	model "Dzaakk/simple-commerce/internal/product/models"
	repo "Dzaakk/simple-commerce/internal/product/repository"
	"fmt"
)

type ProductUseCaseImpl struct {
	repo repo.ProductRepository
}

// Create implements ProductUseCase.
func (p *ProductUseCaseImpl) Create(data model.TProduct) (*int, error) {
	panic("unimplemented")
}

// FilterByPrice implements ProductUseCase.
func (p *ProductUseCaseImpl) FilterByPrice(price int) ([]*model.ProductRes, error) {
	panic("unimplemented")
}

// Update implements ProductUseCase.
func (p *ProductUseCaseImpl) Update(data model.TProduct) error {
	panic("unimplemented")
}

// UpdateStock implements ProductUseCase.
func (p *ProductUseCaseImpl) UpdateStock(id int, stock int) (*int, error) {
	panic("unimplemented")
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
