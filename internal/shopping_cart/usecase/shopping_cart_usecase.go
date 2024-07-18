package shopping_cart

import (
	repoProduct "Dzaakk/synapsis/internal/product/repository"
	model "Dzaakk/synapsis/internal/shopping_cart/models"
	repo "Dzaakk/synapsis/internal/shopping_cart/repository"
	"Dzaakk/synapsis/package/template"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

type ShoppingCartUseCase interface {
	Add(data model.ShoppingCartReq) ([]*model.ShoppingCartItem, error)
	CheckStatus(id, customerId int) (*string, error)
}
type ShoppingCartUseCaseImpl struct {
	repo        repo.ShoppingCartRepository
	repoItem    repo.ShoppingCartItemRepository
	repoProduct repoProduct.ProductRepository
}

func NewShoppingCartUseCase(repo repo.ShoppingCartRepository, repoItem repo.ShoppingCartItemRepository, repoProduct repoProduct.ProductRepository) ShoppingCartUseCase {
	return &ShoppingCartUseCaseImpl{repo, repoItem, repoProduct}
}

func (s *ShoppingCartUseCaseImpl) Add(data model.ShoppingCartReq) ([]*model.ShoppingCartItem, error) {
	customerId, _ := strconv.Atoi(data.CustomerId)
	quantity, _ := strconv.Atoi(data.Quantity)
	productId, _ := strconv.Atoi(data.ProductId)

	//find product and check the stock
	dataProduct, err := s.repoProduct.FindById(productId)
	if err != nil {
		return nil, err
	}
	if quantity > dataProduct.Stock {
		return nil, errors.New("stock product is less than quantity")
	}

	//check is there any chart that active
	shoppingCart, _ := s.repo.FindByCustomerIdAndStatus(customerId, data.Status)
	if shoppingCart == nil {
		newShoppingCart := model.TShoppingCart{
			CustomerId: customerId,
			Status:     data.Status,
			Base: template.Base{
				CreatedBy: fmt.Sprintf("%d", customerId),
				Created:   time.Now(),
			},
		}
		newData, err := s.repo.Create(newShoppingCart)
		if err != nil {
			return nil, err
		}

		cartId, _ := strconv.Atoi(newData.Id)
		// add product to shopping cart item
		newItem := model.TShoppingCartItem{
			ProductId: productId,
			CartId:    cartId,
			Quantity:  quantity,
			Base: template.Base{
				Created:   time.Now(),
				CreatedBy: fmt.Sprintf("%d", customerId),
			},
		}

		//insert into table shopping cart item
		dataItem, err := s.repoItem.Create(newItem)
		if err != nil {
			return nil, err
		}

		item := model.ShoppingCartItem{
			ProductName: dataProduct.ProductName,
			Price:       fmt.Sprintf("%.0f", dataProduct.Price),
			Quantity:    dataItem.Quantity,
		}

		listItem := []*model.ShoppingCartItem{&item}

		return listItem, nil
	} else {
		data.Id = shoppingCart.Id
		if quantity > 0 {
			helperPositiveQuantity(s, data)
		}
		if quantity < 0 {
			helperNegativeQuantity(s, data)
		}

	}

	return nil, nil
}

func helperPositiveQuantity(s *ShoppingCartUseCaseImpl, data model.ShoppingCartReq) ([]*model.ShoppingCartItem, error) {
	customerId, _ := strconv.Atoi(data.CustomerId)
	quantity, _ := strconv.Atoi(data.Quantity)
	productId, _ := strconv.Atoi(data.ProductId)
	cartId, _ := strconv.Atoi(data.Id)
	log.Println("CART ID = ", cartId)
	itemQuantity, _ := s.repoItem.CountQuantityByProductAndCartId(productId, cartId)
	if itemQuantity != 0 {
		quantity += itemQuantity
	}
	if itemQuantity == 0 {
		fmt.Println("QUANTITY = ", quantity)
		newItem := model.TShoppingCartItem{
			ProductId: productId,
			CartId:    cartId,
			Quantity:  quantity,
			Base: template.Base{
				Created:   time.Now(),
				CreatedBy: fmt.Sprintf("%d", customerId),
			},
		}

		fmt.Printf("DATA ITEM =%v\n", newItem)
		//insert into table shopping cart item
		_, err := s.repoItem.Create(newItem)
		if err != nil {
			return nil, err
		}

	}
	fmt.Println("QUANTITY = ", quantity)
	newItem := model.TShoppingCartItem{
		ProductId: productId,
		CartId:    cartId,
		Quantity:  quantity,
	}

	fmt.Printf("DATA ITEM =%v\n", newItem)
	//insert into table shopping cart item
	_, err := s.repoItem.Update(newItem, data.CustomerId)
	if err != nil {
		return nil, err
	}

	var listItem []*model.ShoppingCartItem
	//find list productDetails
	listData, err := s.repoItem.RetrieveCartItemsByCartId(cartId)
	if err != nil {
		return nil, err
	}
	for _, v := range listData {
		item := model.ShoppingCartItem{
			ProductName: v.ProductName,
			Price:       fmt.Sprintf("%0.f", v.Price),
			Quantity:    fmt.Sprintf("%d", v.Quantity),
		}
		listItem = append(listItem, &item)
	}

	return listItem, nil
}

func helperNegativeQuantity(s *ShoppingCartUseCaseImpl, data model.ShoppingCartReq) ([]*model.ShoppingCartItem, error) {
	// customerId, _ := strconv.Atoi(data.CustomerId)
	quantity, _ := strconv.Atoi(data.Quantity)
	productId, _ := strconv.Atoi(data.ProductId)
	cartId, _ := strconv.Atoi(data.Id)

	itemQuantity, _ := s.repoItem.CountQuantityByProductAndCartId(productId, cartId)
	if itemQuantity+quantity < 0 {
		//delete from cart item
		err := s.repoItem.Delete(productId, cartId)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("QUANTITY = ", itemQuantity)
	itemQuantity += quantity
	fmt.Println("QUANTITY 2= ", itemQuantity)
	quantity = itemQuantity
	//find by cart id
	newItem := model.TShoppingCartItem{
		ProductId: productId,
		CartId:    cartId,
		Quantity:  quantity,
	}
	fmt.Printf("DATA ITEM =%v\n", newItem)
	//insert into table shopping cart item
	_, err := s.repoItem.Update(newItem, data.CustomerId)
	if err != nil {
		return nil, err
	}

	var listItem []*model.ShoppingCartItem
	//find list productDetails
	listData, err := s.repoItem.RetrieveCartItemsByCartId(cartId)
	if err != nil {
		return nil, err
	}
	for _, v := range listData {
		if v.Quantity == 0 {
			continue
		}
		item := model.ShoppingCartItem{
			ProductName: v.ProductName,
			Price:       fmt.Sprintf("%0.f", v.Price),
			Quantity:    fmt.Sprintf("%d", v.Quantity),
		}
		listItem = append(listItem, &item)
	}

	return listItem, nil
}

func (s *ShoppingCartUseCaseImpl) CheckStatus(id int, customerId int) (*string, error) {
	status, err := s.repo.CheckStatus(id, customerId)
	if err != nil {
		return nil, err
	}

	return status, nil
}
