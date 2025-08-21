package usecase

import (
	modelProduct "Dzaakk/simple-commerce/internal/product/model"
	"Dzaakk/simple-commerce/internal/shopping_cart/model"
	"Dzaakk/simple-commerce/package/template"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func (s *ShoppingCartUseCaseImpl) parseAndValidateQuantityProduct(ctx context.Context, data model.ShoppingCartReq) (int, int, *modelProduct.TProduct, error) {
	customerID, _ := strconv.Atoi(data.CustomerID)
	quantity, _ := strconv.Atoi(data.Quantity)
	productID, _ := strconv.Atoi(data.ProductID)

	product, err := s.repoProduct.FindByID(ctx, productID)
	if err != nil {
		return 0, 0, nil, err
	}

	if quantity > product.Stock {
		return 0, 0, nil, errors.New("stock product is less than quantity")
	}

	return customerID, quantity, product, nil
}

func (s *ShoppingCartUseCaseImpl) processCartItem(ctx context.Context, cartID, quantity int, product modelProduct.TProduct, data model.ShoppingCartReq) (*model.ShoppingCartItem, error) {

	currentQuantity, _ := s.repoItem.CountQuantityByProductIDAndCartID(ctx, product.ID, cartID)

	newQuantity := currentQuantity + quantity

	if newQuantity == 0 {
		return s.removeItem(ctx, cartID, product)
	}

	if currentQuantity == 0 {
		return s.addNewItem(ctx, cartID, quantity, product, data)
	}

	return s.updateItemQuantity(ctx, cartID, newQuantity, product, data)
}

func (s *ShoppingCartUseCaseImpl) addNewItem(ctx context.Context, cartID, quantity int, product modelProduct.TProduct, data model.ShoppingCartReq) (*model.ShoppingCartItem, error) {
	if _, err := s.CreateCartItem(ctx, data); err != nil {
		return nil, err
	}

	return &model.ShoppingCartItem{
		ProductName: product.ProductName,
		Price:       fmt.Sprintf("%.0f", product.Price),
		Quantity:    fmt.Sprintf("%d", quantity),
	}, nil
}

func (s *ShoppingCartUseCaseImpl) updateItemQuantity(ctx context.Context, cartID, newQuantity int, product modelProduct.TProduct, data model.ShoppingCartReq) (*model.ShoppingCartItem, error) {
	updatedItem := model.TShoppingCartItem{
		ProductID: product.ID,
		CartID:    cartID,
		Quantity:  newQuantity,
		Base: template.Base{
			Updated:   sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedBy: sql.NullString{String: data.CustomerID, Valid: true},
		},
	}

	if _, err := s.repoItem.Update(ctx, updatedItem, data.CustomerID); err != nil {
		return nil, err
	}

	return &model.ShoppingCartItem{
		ProductName: product.ProductName,
		Price:       fmt.Sprintf("%.0f", product.Price),
		Quantity:    fmt.Sprintf("%d", newQuantity),
	}, nil

}

func (s *ShoppingCartUseCaseImpl) removeItem(ctx context.Context, cartID int, product modelProduct.TProduct) (*model.ShoppingCartItem, error) {
	if err := s.repoItem.Delete(ctx, product.ID, cartID); err != nil {
		return nil, err
	}

	return nil, nil
}
