package service

import (
	"context"

	"Dzaakk/simple-commerce/internal/catalog/dto"
)

type CategoryServiceImpl struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryServiceImpl {
	return &CategoryServiceImpl{repo: repo}
}

func (s *CategoryServiceImpl) FindAll(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		result = append(result, dto.ToCategoryResponse(category))
	}
	return result, nil
}
