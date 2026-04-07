package service

import (
	"Dzaakk/simple-commerce/internal/cart/dto"
	cartModel "Dzaakk/simple-commerce/internal/cart/model"
	"context"
	"errors"
)

type CartServiceImpl struct {
	CartRepo      CartRepository
	CartItemRepo  CartItemRepository
	ProductSvc    ProductService
	InventorySvc  InventoryService
}

func NewCartService(cartRepo CartRepository, cartItemRepo CartItemRepository, productSvc ProductService, inventorySvc InventoryService) CartService {
	return &CartServiceImpl{
		CartRepo:     cartRepo,
		CartItemRepo: cartItemRepo,
		ProductSvc:   productSvc,
		InventorySvc: inventorySvc,
	}
}

func (s *CartServiceImpl) GetCartItems(ctx context.Context, customerID string) (*dto.CartRes, error) {
	if customerID == "" {
		return nil, errors.New("invalid parameter customer id")
	}

	cart, err := s.CartRepo.GetCartByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return &dto.CartRes{CustomerID: customerID, Items: []dto.CartItemRes{}}, nil
	}

	items, err := s.CartItemRepo.GetCartItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	return buildCartRes(cart, items), nil
}

func (s *CartServiceImpl) AddItem(ctx context.Context, customerID string, productID string, quantity int) (*dto.CartRes, error) {
	if customerID == "" {
		return nil, errors.New("invalid parameter customer id")
	}
	if productID == "" {
		return nil, errors.New("invalid parameter product id")
	}
	if quantity <= 0 {
		return nil, errors.New("invalid parameter quantity")
	}

	cart, err := s.CartRepo.GetOrCreateCart(ctx, customerID)
	if err != nil {
		return nil, err
	}

	existingItems, err := s.CartItemRepo.GetCartItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	existingItem := findCartItem(existingItems, productID)
	newQty := quantity
	priceSnapshot := 0.0
	hasSnapshot := false
	if existingItem != nil {
		newQty = existingItem.Quantity + quantity
		priceSnapshot = existingItem.PriceSnapshot
		hasSnapshot = true
	}

	if err := s.ensureStock(ctx, productID, newQty); err != nil {
		return nil, err
	}

	if !hasSnapshot {
		product, err := s.ProductSvc.FindByID(ctx, productID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, errors.New("product not found")
		}
		priceSnapshot = product.Price
	}

	if err := s.CartItemRepo.UpsertItem(ctx, cart.ID, productID, newQty, priceSnapshot); err != nil {
		return nil, err
	}

	items, err := s.CartItemRepo.GetCartItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	return buildCartRes(cart, items), nil
}

func (s *CartServiceImpl) UpdateItem(ctx context.Context, customerID string, productID string, quantity int) (*dto.CartRes, error) {
	if customerID == "" {
		return nil, errors.New("invalid parameter customer id")
	}
	if productID == "" {
		return nil, errors.New("invalid parameter product id")
	}
	if quantity <= 0 {
		return nil, errors.New("invalid parameter quantity")
	}

	cart, err := s.CartRepo.GetCartByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, errors.New("cart not found")
	}

	items, err := s.CartItemRepo.GetCartItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	existingItem := findCartItem(items, productID)
	if existingItem == nil {
		return nil, errors.New("cart item not found")
	}

	if err := s.ensureStock(ctx, productID, quantity); err != nil {
		return nil, err
	}

	if err := s.CartItemRepo.UpsertItem(ctx, cart.ID, productID, quantity, existingItem.PriceSnapshot); err != nil {
		return nil, err
	}

	items, err = s.CartItemRepo.GetCartItems(ctx, cart.ID)
	if err != nil {
		return nil, err
	}

	return buildCartRes(cart, items), nil
}

func (s *CartServiceImpl) DeleteItem(ctx context.Context, customerID string, productID string) error {
	if customerID == "" {
		return errors.New("invalid parameter customer id")
	}
	if productID == "" {
		return errors.New("invalid parameter product id")
	}

	cart, err := s.CartRepo.GetCartByCustomerID(ctx, customerID)
	if err != nil {
		return err
	}
	if cart == nil {
		return errors.New("cart not found")
	}

	return s.CartItemRepo.DeleteItem(ctx, cart.ID, productID)
}

func (s *CartServiceImpl) ClearItems(ctx context.Context, customerID string) error {
	if customerID == "" {
		return errors.New("invalid parameter customer id")
	}

	cart, err := s.CartRepo.GetCartByCustomerID(ctx, customerID)
	if err != nil {
		return err
	}
	if cart == nil {
		return errors.New("cart not found")
	}

	return s.CartItemRepo.ClearItems(ctx, cart.ID)
}

func (s *CartServiceImpl) ensureStock(ctx context.Context, productID string, quantity int) error {
	inventory, err := s.InventorySvc.FindByProductID(ctx, productID)
	if err != nil {
		return err
	}
	if inventory == nil {
		return errors.New("inventory not found")
	}

	available := inventory.StockQuantity - inventory.ReservedQuantity
	if available < quantity {
		return errors.New("stock product is less than quantity")
	}

	return nil
}

func buildCartRes(cart *cartModel.Cart, items []*cartModel.CartItem) *dto.CartRes {
	res := &dto.CartRes{
		CartID:     cart.ID,
		CustomerID: cart.CustomerID,
		Items:      make([]dto.CartItemRes, 0, len(items)),
	}

	var total float64
	for _, item := range items {
		if item == nil {
			continue
		}
		subtotal := item.PriceSnapshot * float64(item.Quantity)
		res.Items = append(res.Items, dto.CartItemRes{
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			PriceSnapshot: item.PriceSnapshot,
			Subtotal:      subtotal,
		})
		total += subtotal
	}

	res.Total = total
	return res
}

func findCartItem(items []*cartModel.CartItem, productID string) *cartModel.CartItem {
	for _, item := range items {
		if item == nil {
			continue
		}
		if item.ProductID == productID {
			return item
		}
	}

	return nil
}
