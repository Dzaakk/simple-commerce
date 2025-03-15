package usecases

import (
	modelProduct "Dzaakk/simple-commerce/internal/product/models"
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	template "Dzaakk/simple-commerce/package/templates"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func (s *ShoppingCartUseCaseImpl) parseAndValidateQuantityProduct(ctx context.Context, data model.ShoppingCartReq) (int, int, *modelProduct.TProduct, error) {
	customerId, _ := strconv.Atoi(data.CustomerId)
	quantity, _ := strconv.Atoi(data.Quantity)
	productId, _ := strconv.Atoi(data.ProductId)

	product, err := s.repoProduct.FindById(ctx, productId)
	if err != nil {
		return 0, 0, nil, err
	}

	if quantity > product.Stock {
		return 0, 0, nil, errors.New("stock product is less than quantity")
	}

	return customerId, quantity, product, nil
}

func (s *ShoppingCartUseCaseImpl) processCartItem(ctx context.Context, cartId, quantity int, product modelProduct.TProduct, data model.ShoppingCartReq) (*model.ShoppingCartItem, error) {

	currentQuantity, _ := s.repoItem.CountQuantityByProductAndCartId(ctx, product.Id, cartId)

	newQuantity := currentQuantity + quantity

	if newQuantity == 0 {
		return s.removeItem(ctx, cartId, product)
	}

	if currentQuantity == 0 {
		return s.addNewItem(ctx, cartId, quantity, product, data)
	}

	return s.updateItemQuantity(ctx, cartId, newQuantity, product, data)
}

func (s *ShoppingCartUseCaseImpl) addNewItem(ctx context.Context, cartId, quantity int, product modelProduct.TProduct, data model.ShoppingCartReq) (*model.ShoppingCartItem, error) {
	if _, err := s.CreateCartItem(ctx, data); err != nil {
		return nil, err
	}

	return &model.ShoppingCartItem{
		ProductName: product.ProductName,
		Price:       fmt.Sprintf("%.0f", product.Price),
		Quantity:    fmt.Sprintf("%d", quantity),
	}, nil
}

func (s *ShoppingCartUseCaseImpl) updateItemQuantity(ctx context.Context, cartId, newQuantity int, product modelProduct.TProduct, data model.ShoppingCartReq) (*model.ShoppingCartItem, error) {
	updatedItem := model.TShoppingCartItem{
		ProductId: product.Id,
		CartId:    cartId,
		Quantity:  newQuantity,
		Base: template.Base{
			Updated:   sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedBy: sql.NullString{String: data.CustomerId, Valid: true},
		},
	}

	if _, err := s.repoItem.Update(ctx, updatedItem, data.CustomerId); err != nil {
		return nil, err
	}

	return &model.ShoppingCartItem{
		ProductName: product.ProductName,
		Price:       fmt.Sprintf("%.0f", product.Price),
		Quantity:    fmt.Sprintf("%d", newQuantity),
	}, nil

}

func (s *ShoppingCartUseCaseImpl) removeItem(ctx context.Context, cartId int, product modelProduct.TProduct) (*model.ShoppingCartItem, error) {
	if err := s.repoItem.Delete(ctx, product.Id, cartId); err != nil {
		return nil, err
	}

	return nil, nil
}
