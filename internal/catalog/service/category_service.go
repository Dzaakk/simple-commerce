package service

import (
	"Dzaakk/simple-commerce/internal/catalog/dto"
	"Dzaakk/simple-commerce/internal/catalog/model"
	"Dzaakk/simple-commerce/package/response"
	"context"
	"net/http"
	"sort"
)

type CategoryServiceImpl struct {
	Repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &CategoryServiceImpl{Repo: repo}
}

func (c *CategoryServiceImpl) Create(ctx context.Context, req *dto.CreateCategoryReq) (int64, error) {
	data := req.ToCreateData()

	id, err := c.Repo.Create(ctx, data)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *CategoryServiceImpl) FindAll(ctx context.Context) ([]*dto.CategoryTree, error) {
	data, err := c.Repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return buildCategoryTrees(data)
}

func (c *CategoryServiceImpl) FindByID(ctx context.Context, categoryID int64) (*dto.CategoryTree, error) {
	if categoryID <= 0 {
		return nil, response.NewAppError(http.StatusBadRequest, "invalid parameter category id")
	}

	data, err := c.Repo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, response.NewAppError(http.StatusNotFound, "category not found")
	}

	category := dto.CategoryTree{
		ID:       data.ID,
		ParentID: data.ParentID,
		Name:     data.Name,
		Slug:     data.Slug,
		Depth:    0,
	}

	return &category, nil
}

func buildCategoryTrees(categories []*model.Category) ([]*dto.CategoryTree, error) {
	if len(categories) == 0 {
		return []*dto.CategoryTree{}, nil
	}

	children := map[int64][]*model.Category{}
	roots := make([]*model.Category, 0)

	for _, category := range categories {
		if category == nil {
			continue
		}
		if category.ParentID == nil {
			roots = append(roots, category)
			continue
		}
		children[*category.ParentID] = append(children[*category.ParentID], category)
	}

	result := make([]*dto.CategoryTree, 0, len(categories))
	visited := map[int64]bool{}
	onStack := map[int64]bool{}

	var walk func(node *model.Category, depth int) error
	walk = func(node *model.Category, depth int) error {
		if node == nil {
			return nil
		}
		if onStack[node.ID] {
			return response.NewAppError(http.StatusBadRequest, "category tree has a cycle")
		}
		if visited[node.ID] {
			return nil
		}

		visited[node.ID] = true
		onStack[node.ID] = true

		result = append(result, &dto.CategoryTree{
			ID:       node.ID,
			ParentID: node.ParentID,
			Name:     node.Name,
			Slug:     node.Slug,
			Depth:    depth,
		})

		for _, child := range children[node.ID] {
			if err := walk(child, depth+1); err != nil {
				return err
			}
		}

		onStack[node.ID] = false

		return nil
	}

	for _, root := range roots {
		if err := walk(root, 0); err != nil {
			return nil, err
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Depth == result[j].Depth {
			return result[i].ID < result[j].ID
		}
		return result[i].Depth < result[j].Depth
	})

	return result, nil
}
