package usecase

import (
	repoProduct "Dzaakk/simple-commerce/internal/product/repositories"
	model "Dzaakk/simple-commerce/internal/shopping_cart/models"
	repo "Dzaakk/simple-commerce/internal/shopping_cart/repositories"
	template "Dzaakk/simple-commerce/package/templates"
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

func (s *ShoppingCartUseCaseImpl) DeleteShoppingList(data model.DeleteReq) error {
	cartId, _ := strconv.Atoi(data.CartId)

	err := s.repoItem.DeleteAll(cartId)
	if err != nil {
		return err
	}

	_, err = s.repo.UpdateStatusById(cartId, "Inactive", data.CustomerId)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShoppingCartUseCaseImpl) GetListItem(customerId int) ([]*model.ListCartItemRes, error) {

	cart, _ := s.repo.FindByCustomerIdAndStatus(customerId, "Active")
	if cart == nil {
		return nil, errors.New("cart is empty")
	}

	listData, err := s.repoItem.RetrieveCartItemsByCartId(cart.Id)
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

func (s *ShoppingCartUseCaseImpl) Add(data model.ShoppingCartReq) (*model.ShoppingCartItem, error) {
	customerId, _ := strconv.Atoi(data.CustomerId)
	quantity, _ := strconv.Atoi(data.Quantity)
	productId, _ := strconv.Atoi(data.ProductId)

	dataProduct, err := s.repoProduct.FindById(productId) //find product and check the stock
	if err != nil {
		return nil, err
	}
	if quantity > dataProduct.Stock {
		return nil, errors.New("stock product is less than quantity")
	}

	shoppingCart, _ := s.repo.FindByCustomerIdAndStatus(customerId, "Active") //check is there any chart that active
	if shoppingCart == nil {
		cart, err := s.CreateCart(customerId) // create new cart
		if err != nil {
			return nil, err
		}
		data.Id = strconv.Itoa(cart.Id)
		_, err = s.CreateCartItem(data) // insert product to cart item
		if err != nil {
			return nil, err
		}

		item := model.ShoppingCartItem{
			ProductName: dataProduct.ProductName,
			Price:       fmt.Sprintf("%.0f", dataProduct.Price),
			Quantity:    fmt.Sprintf("%d", quantity),
			NewCartId:   fmt.Sprintf("%d", cart.Id),
		}

		return &item, nil
	} else {
		itemQuantity, _ := s.repoItem.CountQuantityByProductAndCartId(productId, shoppingCart.Id) //check current quantity product
		if itemQuantity+quantity == 0 {
			err = s.repoItem.Delete(productId, shoppingCart.Id) //delete from cart item
			if err != nil {
				return nil, err
			}

			count, err := s.repoItem.CountByCartId(shoppingCart.Id) //count cart item base on chart_id
			if err != nil {
				return nil, err
			}
			if count == 0 {
				_, err = s.repo.UpdateStatusById(shoppingCart.Id, "Inactive", data.CustomerId) //update status shopping chart to inactive
				if err != nil {
					return nil, err
				}
			}

			return nil, nil

		} else if itemQuantity == 0 { //create new if the product is not on the cart item
			data.Id = fmt.Sprintf("%d", shoppingCart.Id)
			_, err = s.CreateCartItem(data)
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
				ProductId: productId,
				CartId:    shoppingCart.Id,
				Quantity:  itemQuantity,
				Base: template.Base{
					Updated:   sql.NullTime{Time: time.Now(), Valid: true},
					UpdatedBy: sql.NullString{String: data.CustomerId, Valid: true},
				},
			}

			_, err = s.repoItem.Update(newItem, data.CustomerId) // update quantity product at shopping cart item
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

func (s *ShoppingCartUseCaseImpl) CreateCart(customerId int) (*model.TShoppingCart, error) {
	newData := model.TShoppingCart{
		CustomerId: customerId,
		Status:     "Active",
		Base: template.Base{
			CreatedBy: fmt.Sprintf("%d", customerId),
			Created:   time.Now(),
		},
	}
	data, err := s.repo.Create(newData)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ShoppingCartUseCaseImpl) CreateCartItem(data model.ShoppingCartReq) (*model.TShoppingCartItem, error) {
	customerId, _ := strconv.Atoi(data.CustomerId)
	quantity, _ := strconv.Atoi(data.Quantity)
	productId, _ := strconv.Atoi(data.ProductId)
	cartId, _ := strconv.Atoi(data.Id)

	newItem := model.TShoppingCartItem{
		ProductId: productId,
		CartId:    cartId,
		Quantity:  quantity,
		Base: template.Base{
			Created:   time.Now(),
			CreatedBy: fmt.Sprintf("%d", customerId),
		},
	}

	cartItem, err := s.repoItem.Create(newItem) //insert into table shopping cart item
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (s *ShoppingCartUseCaseImpl) CheckStatus(id int, customerId int) (string, error) {
	status, err := s.repo.CheckStatus(id, customerId)
	if err != nil {
		return "", err
	}

	return status, nil
}
