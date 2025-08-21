package usecase

import (
	repoProduct "Dzaakk/simple-commerce/internal/product/repository"
	"Dzaakk/simple-commerce/internal/shopping_cart/model"
	repo "Dzaakk/simple-commerce/internal/shopping_cart/repository"
	"Dzaakk/simple-commerce/package/template"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type ShoppingCartUseCaseImpl struct {
	repo        repo.ShoppingCartRepository
	repoItem    repo.ShoppingCartItemRepository
	repoProduct repoProduct.ProductRepository
}

func NewShoppingCartUseCase(repo repo.ShoppingCartRepository, repoItem repo.ShoppingCartItemRepository, repoProduct repoProduct.ProductRepository) ShoppingCartUseCase {
	return &ShoppingCartUseCaseImpl{repo, repoItem, repoProduct}
}

func (s *ShoppingCartUseCaseImpl) DeleteShoppingList(ctx context.Context, data model.DeleteReq) error {
	cartID, _ := strconv.Atoi(data.CartID)

	err := s.repoItem.DeleteAll(ctx, cartID)
	if err != nil {
		return err
	}

	_, err = s.repo.UpdateStatusByCartID(ctx, cartID, "Inactive", data.CustomerID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShoppingCartUseCaseImpl) GetListItem(ctx context.Context, customerID int) ([]*model.ListCartItemRes, error) {

	cart, _ := s.repo.FindByStatusAndCustomerID(ctx, customerID, "Active")
	if cart == nil {
		return nil, errors.New("cart is empty")
	}

	listData, err := s.repoItem.RetrieveCartItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	var listItem []*model.ListCartItemRes

	for _, d := range listData {
		dataProduct := model.ShoppingCartItem{
			ProductName: d.ProductName,
			Price:       fmt.Sprintf("%.0f", d.Price),
			Quantity:    fmt.Sprintf("%d", d.Quantity),
		}
		totalPrice := d.Quantity * int(d.Price)

		item := model.ListCartItemRes{
			Product:    dataProduct,
			TotalPrice: fmt.Sprintf("%d", totalPrice),
		}

		listItem = append(listItem, &item)
	}

	return listItem, nil
}

func (s *ShoppingCartUseCaseImpl) AddV2(ctx context.Context, data model.ShoppingCartReq) (*model.ShoppingCartItem, error) {
	customerID, quantity, product, err := s.parseAndValidateQuantityProduct(ctx, data)
	if err != nil {
		return nil, err
	}

	cart, err := s.repo.FindByCustomerID(ctx, customerID) //find customer cart id
	if err != nil {
		return nil, errors.New("failed retrieve customer cart")
	}

	return s.processCartItem(ctx, cart.ID, quantity, *product, data)
}

func (s *ShoppingCartUseCaseImpl) Add(ctx context.Context, data model.ShoppingCartReq) (*model.ShoppingCartItem, error) {
	customerID, _ := strconv.Atoi(data.CustomerID)
	quantity, _ := strconv.Atoi(data.Quantity)
	productID, _ := strconv.Atoi(data.ProductID)

	dataProduct, err := s.repoProduct.FindByID(ctx, productID) //find product and check the stock
	if err != nil {
		return nil, err
	}
	if quantity > dataProduct.Stock {
		return nil, errors.New("stock product is less than quantity")
	}

	shoppingCart, _ := s.repo.FindByStatusAndCustomerID(ctx, customerID, "Active") //check is there any chart that active
	if shoppingCart == nil {
		cart, err := s.CreateCart(ctx, customerID) // create new cart
		if err != nil {
			return nil, err
		}
		data.ShoppingCartID = strconv.Itoa(cart.ID)
		_, err = s.CreateCartItem(ctx, data) // insert product to cart item
		if err != nil {
			return nil, err
		}

		item := model.ShoppingCartItem{
			ProductName: dataProduct.ProductName,
			Price:       fmt.Sprintf("%.0f", dataProduct.Price),
			Quantity:    fmt.Sprintf("%d", quantity),
			NewCartID:   fmt.Sprintf("%d", cart.ID),
		}

		return &item, nil
	} else {
		itemQuantity, _ := s.repoItem.CountQuantityByProductIDAndCartID(ctx, productID, shoppingCart.ID) //check current quantity product
		if itemQuantity+quantity == 0 {
			err = s.repoItem.Delete(ctx, productID, shoppingCart.ID) //delete from cart item
			if err != nil {
				return nil, err
			}

			count, err := s.repoItem.CountByCartID(ctx, shoppingCart.ID) //count cart item base on chart_id
			if err != nil {
				return nil, err
			}
			if count == 0 {
				_, err = s.repo.UpdateStatusByCartID(ctx, shoppingCart.ID, "Inactive", data.CustomerID) //update status shopping chart to inactive
				if err != nil {
					return nil, err
				}
			}

			return nil, nil

		} else if itemQuantity == 0 { //create new if the product is not on the cart item
			data.ShoppingCartID = fmt.Sprintf("%d", shoppingCart.ID)
			_, err = s.CreateCartItem(ctx, data)
			if err != nil {
				return nil, err
			}

			item := model.ShoppingCartItem{
				ProductName: dataProduct.ProductName,
				Price:       fmt.Sprintf("%.0f", dataProduct.Price),
				Quantity:    fmt.Sprintf("%d", quantity),
			}

			return &item, nil
		} else { //update quantity if the product is exist on the cart item
			itemQuantity += quantity
			newItem := model.TShoppingCartItem{
				ProductID: productID,
				CartID:    shoppingCart.ID,
				Quantity:  itemQuantity,
				Base: template.Base{
					Updated:   sql.NullTime{Time: time.Now(), Valid: true},
					UpdatedBy: sql.NullString{String: data.CustomerID, Valid: true},
				},
			}

			_, err = s.repoItem.Update(ctx, newItem, data.CustomerID) // update quantity product at shopping cart item
			if err != nil {
				return nil, err
			}
			item := model.ShoppingCartItem{
				ProductName: dataProduct.ProductName,
				Price:       fmt.Sprintf("%.0f", dataProduct.Price),
				Quantity:    fmt.Sprintf("%d", itemQuantity),
			}

			return &item, nil
		}

	}
}

func (s *ShoppingCartUseCaseImpl) CreateCart(ctx context.Context, customerID int) (*model.TShoppingCart, error) {
	newData := model.TShoppingCart{
		CustomerID: customerID,
		Status:     "Active",
		Base: template.Base{
			CreatedBy: fmt.Sprintf("%d", customerID),
			Created:   time.Now(),
		},
	}
	data, err := s.repo.Create(ctx, newData)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ShoppingCartUseCaseImpl) CreateCartItem(ctx context.Context, data model.ShoppingCartReq) (*model.TShoppingCartItem, error) {
	customerID, _ := strconv.Atoi(data.CustomerID)
	quantity, _ := strconv.Atoi(data.Quantity)
	productID, _ := strconv.Atoi(data.ProductID)
	cartID, _ := strconv.Atoi(data.ShoppingCartID)

	newItem := model.TShoppingCartItem{
		ProductID: productID,
		CartID:    cartID,
		Quantity:  quantity,
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: fmt.Sprintf("%d", customerID),
		},
	}

	cartItem, err := s.repoItem.Create(ctx, newItem) //insert into table shopping cart item
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (s *ShoppingCartUseCaseImpl) CheckStatus(ctx context.Context, cartID int, customerID int) (string, error) {
	status, err := s.repo.CheckStatus(ctx, cartID, customerID)
	if err != nil {
		return "", err
	}

	return status, nil
}
