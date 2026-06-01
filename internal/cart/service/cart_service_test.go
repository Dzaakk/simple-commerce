package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"Dzaakk/simple-commerce/internal/cart/dto"
	cartModel "Dzaakk/simple-commerce/internal/cart/model"
	catalogDto "Dzaakk/simple-commerce/internal/catalog/dto"
	catalogModel "Dzaakk/simple-commerce/internal/catalog/model"
	"Dzaakk/simple-commerce/package/response"
)

type mockCartRepository struct {
	getCartByCustomerIDFn func(context.Context, string) (*cartModel.Cart, error)
	getOrCreateCartFn     func(context.Context, string) (*cartModel.Cart, error)
}

func (m *mockCartRepository) GetCartByCustomerID(ctx context.Context, customerID string) (*cartModel.Cart, error) {
	if m.getCartByCustomerIDFn == nil {
		return nil, errors.New("unexpected GetCartByCustomerID call")
	}
	return m.getCartByCustomerIDFn(ctx, customerID)
}

func (m *mockCartRepository) GetOrCreateCart(ctx context.Context, customerID string) (*cartModel.Cart, error) {
	if m.getOrCreateCartFn == nil {
		return nil, errors.New("unexpected GetOrCreateCart call")
	}
	return m.getOrCreateCartFn(ctx, customerID)
}

type mockCartItemRepository struct {
	getCartItemsFn func(context.Context, string, int) ([]*cartModel.CartItem, error)
	upsertItemFn   func(context.Context, string, string, int, float64) error
	deleteItemFn   func(context.Context, string, string) error
	clearItemsFn   func(context.Context, string) error

	getCartItemsCalls int
}

func (m *mockCartItemRepository) GetCartItems(ctx context.Context, cartID string) ([]*cartModel.CartItem, error) {
	if m.getCartItemsFn == nil {
		return nil, errors.New("unexpected GetCartItems call")
	}
	m.getCartItemsCalls++
	return m.getCartItemsFn(ctx, cartID, m.getCartItemsCalls)
}

func (m *mockCartItemRepository) UpsertItem(ctx context.Context, cartID string, productID string, quantity int, priceSnapshot float64) error {
	if m.upsertItemFn == nil {
		return errors.New("unexpected UpsertItem call")
	}
	return m.upsertItemFn(ctx, cartID, productID, quantity, priceSnapshot)
}

func (m *mockCartItemRepository) DeleteItem(ctx context.Context, cartID string, productID string) error {
	if m.deleteItemFn == nil {
		return errors.New("unexpected DeleteItem call")
	}
	return m.deleteItemFn(ctx, cartID, productID)
}

func (m *mockCartItemRepository) ClearItems(ctx context.Context, cartID string) error {
	if m.clearItemsFn == nil {
		return errors.New("unexpected ClearItems call")
	}
	return m.clearItemsFn(ctx, cartID)
}

type mockProductService struct {
	findByIDFn func(context.Context, string) (*catalogDto.ProductRes, error)
}

func (m *mockProductService) FindByID(ctx context.Context, productID string) (*catalogDto.ProductRes, error) {
	if m.findByIDFn == nil {
		return nil, errors.New("unexpected Product FindByID call")
	}
	return m.findByIDFn(ctx, productID)
}

type mockInventoryService struct {
	findByProductIDFn func(context.Context, string) (*catalogModel.Inventory, error)
}

func (m *mockInventoryService) FindByProductID(ctx context.Context, productID string) (*catalogModel.Inventory, error) {
	if m.findByProductIDFn == nil {
		return nil, errors.New("unexpected Inventory FindByProductID call")
	}
	return m.findByProductIDFn(ctx, productID)
}

func TestCartServiceGetCartItemsRejectsInvalidCustomerID(t *testing.T) {
	service := NewCartService(&mockCartRepository{}, &mockCartItemRepository{}, &mockProductService{}, &mockInventoryService{})

	got, err := service.GetCartItems(context.Background(), "")
	assertCartAppError(t, err, http.StatusBadRequest, "invalid parameter customer id")
	if got != nil {
		t.Fatalf("cart = %#v, want nil", got)
	}
}

func TestCartServiceGetCartItemsReturnsEmptyCartWhenCartNotFound(t *testing.T) {
	repo := &mockCartRepository{
		getCartByCustomerIDFn: func(_ context.Context, customerID string) (*cartModel.Cart, error) {
			if customerID != "customer-1" {
				t.Fatalf("customer id = %q, want %q", customerID, "customer-1")
			}
			return nil, nil
		},
	}

	got, err := NewCartService(repo, &mockCartItemRepository{}, &mockProductService{}, &mockInventoryService{}).
		GetCartItems(context.Background(), "customer-1")
	if err != nil {
		t.Fatalf("GetCartItems returned error: %v", err)
	}
	if got.CustomerID != "customer-1" {
		t.Fatalf("customer id = %q, want %q", got.CustomerID, "customer-1")
	}
	if got.Items == nil {
		t.Fatal("items must not be nil")
	}
	if len(got.Items) != 0 {
		t.Fatalf("items count = %d, want 0", len(got.Items))
	}
	if got.Total != 0 {
		t.Fatalf("total = %.2f, want 0", got.Total)
	}
}

func TestCartServiceGetCartItemsBuildsResponse(t *testing.T) {
	cart := &cartModel.Cart{ID: "cart-1", CustomerID: "customer-1"}
	cartRepo := &mockCartRepository{
		getCartByCustomerIDFn: func(context.Context, string) (*cartModel.Cart, error) {
			return cart, nil
		},
	}
	itemRepo := &mockCartItemRepository{
		getCartItemsFn: func(_ context.Context, cartID string, call int) ([]*cartModel.CartItem, error) {
			if call != 1 {
				t.Fatalf("GetCartItems call = %d, want 1", call)
			}
			if cartID != cart.ID {
				t.Fatalf("cart id = %q, want %q", cartID, cart.ID)
			}
			return []*cartModel.CartItem{
				{ProductID: "product-1", Quantity: 2, PriceSnapshot: 10000},
				nil,
				{ProductID: "product-2", Quantity: 1, PriceSnapshot: 25000},
			}, nil
		},
	}

	got, err := NewCartService(cartRepo, itemRepo, &mockProductService{}, &mockInventoryService{}).
		GetCartItems(context.Background(), "customer-1")
	if err != nil {
		t.Fatalf("GetCartItems returned error: %v", err)
	}
	assertCartResponse(t, got, "cart-1", "customer-1", 45000, []dto.CartItemRes{
		{ProductID: "product-1", Quantity: 2, PriceSnapshot: 10000, Subtotal: 20000},
		{ProductID: "product-2", Quantity: 1, PriceSnapshot: 25000, Subtotal: 25000},
	})
}

func TestCartServiceAddItemRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name       string
		customerID string
		productID  string
		quantity   int
		message    string
	}{
		{name: "empty customer id", productID: "product-1", quantity: 1, message: "invalid parameter customer id"},
		{name: "empty product id", customerID: "customer-1", quantity: 1, message: "invalid parameter product id"},
		{name: "zero quantity", customerID: "customer-1", productID: "product-1", message: "invalid parameter quantity"},
		{name: "negative quantity", customerID: "customer-1", productID: "product-1", quantity: -1, message: "invalid parameter quantity"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewCartService(&mockCartRepository{}, &mockCartItemRepository{}, &mockProductService{}, &mockInventoryService{})

			got, err := service.AddItem(context.Background(), tt.customerID, tt.productID, tt.quantity)
			assertCartAppError(t, err, http.StatusBadRequest, tt.message)
			if got != nil {
				t.Fatalf("cart = %#v, want nil", got)
			}
		})
	}
}

func TestCartServiceAddItemCreatesNewItemWithProductPrice(t *testing.T) {
	cart := &cartModel.Cart{ID: "cart-1", CustomerID: "customer-1"}
	cartRepo := &mockCartRepository{
		getOrCreateCartFn: func(_ context.Context, customerID string) (*cartModel.Cart, error) {
			if customerID != cart.CustomerID {
				t.Fatalf("customer id = %q, want %q", customerID, cart.CustomerID)
			}
			return cart, nil
		},
	}
	itemRepo := &mockCartItemRepository{
		getCartItemsFn: func(_ context.Context, cartID string, call int) ([]*cartModel.CartItem, error) {
			if cartID != cart.ID {
				t.Fatalf("cart id = %q, want %q", cartID, cart.ID)
			}
			if call == 1 {
				return []*cartModel.CartItem{}, nil
			}
			return []*cartModel.CartItem{{ProductID: "product-1", Quantity: 2, PriceSnapshot: 15000}}, nil
		},
		upsertItemFn: func(_ context.Context, cartID string, productID string, quantity int, priceSnapshot float64) error {
			if cartID != cart.ID || productID != "product-1" || quantity != 2 || priceSnapshot != 15000 {
				t.Fatalf("upsert args = %q %q %d %.2f", cartID, productID, quantity, priceSnapshot)
			}
			return nil
		},
	}
	productSvc := &mockProductService{
		findByIDFn: func(_ context.Context, productID string) (*catalogDto.ProductRes, error) {
			if productID != "product-1" {
				t.Fatalf("product id = %q, want product-1", productID)
			}
			return &catalogDto.ProductRes{ID: productID, Price: 15000}, nil
		},
	}
	inventorySvc := &mockInventoryService{
		findByProductIDFn: func(_ context.Context, productID string) (*catalogModel.Inventory, error) {
			if productID != "product-1" {
				t.Fatalf("product id = %q, want product-1", productID)
			}
			return &catalogModel.Inventory{ProductID: productID, StockQuantity: 10, ReservedQuantity: 3}, nil
		},
	}

	got, err := NewCartService(cartRepo, itemRepo, productSvc, inventorySvc).
		AddItem(context.Background(), "customer-1", "product-1", 2)
	if err != nil {
		t.Fatalf("AddItem returned error: %v", err)
	}
	assertCartResponse(t, got, "cart-1", "customer-1", 30000, []dto.CartItemRes{
		{ProductID: "product-1", Quantity: 2, PriceSnapshot: 15000, Subtotal: 30000},
	})
}

func TestCartServiceAddItemIncrementsExistingItemAndReusesPriceSnapshot(t *testing.T) {
	cart := &cartModel.Cart{ID: "cart-1", CustomerID: "customer-1"}
	productLookupCalled := false
	cartRepo := &mockCartRepository{
		getOrCreateCartFn: func(context.Context, string) (*cartModel.Cart, error) {
			return cart, nil
		},
	}
	itemRepo := &mockCartItemRepository{
		getCartItemsFn: func(_ context.Context, _ string, call int) ([]*cartModel.CartItem, error) {
			if call == 1 {
				return []*cartModel.CartItem{{ProductID: "product-1", Quantity: 3, PriceSnapshot: 9000}}, nil
			}
			return []*cartModel.CartItem{{ProductID: "product-1", Quantity: 5, PriceSnapshot: 9000}}, nil
		},
		upsertItemFn: func(_ context.Context, _ string, productID string, quantity int, priceSnapshot float64) error {
			if productID != "product-1" || quantity != 5 || priceSnapshot != 9000 {
				t.Fatalf("upsert args = %q %d %.2f, want product-1 5 9000", productID, quantity, priceSnapshot)
			}
			return nil
		},
	}
	productSvc := &mockProductService{
		findByIDFn: func(context.Context, string) (*catalogDto.ProductRes, error) {
			productLookupCalled = true
			return nil, nil
		},
	}
	inventorySvc := &mockInventoryService{
		findByProductIDFn: func(context.Context, string) (*catalogModel.Inventory, error) {
			return &catalogModel.Inventory{StockQuantity: 8, ReservedQuantity: 1}, nil
		},
	}

	got, err := NewCartService(cartRepo, itemRepo, productSvc, inventorySvc).
		AddItem(context.Background(), "customer-1", "product-1", 2)
	if err != nil {
		t.Fatalf("AddItem returned error: %v", err)
	}
	if productLookupCalled {
		t.Fatal("product lookup must not be called when existing item has price snapshot")
	}
	assertCartResponse(t, got, "cart-1", "customer-1", 45000, []dto.CartItemRes{
		{ProductID: "product-1", Quantity: 5, PriceSnapshot: 9000, Subtotal: 45000},
	})
}

func TestCartServiceAddItemReturnsConflictWhenStockInsufficient(t *testing.T) {
	cart := &cartModel.Cart{ID: "cart-1", CustomerID: "customer-1"}
	service := NewCartService(
		&mockCartRepository{
			getOrCreateCartFn: func(context.Context, string) (*cartModel.Cart, error) {
				return cart, nil
			},
		},
		&mockCartItemRepository{
			getCartItemsFn: func(context.Context, string, int) ([]*cartModel.CartItem, error) {
				return []*cartModel.CartItem{}, nil
			},
		},
		&mockProductService{},
		&mockInventoryService{
			findByProductIDFn: func(context.Context, string) (*catalogModel.Inventory, error) {
				return &catalogModel.Inventory{StockQuantity: 3, ReservedQuantity: 1}, nil
			},
		},
	)

	got, err := service.AddItem(context.Background(), "customer-1", "product-1", 3)
	assertCartAppError(t, err, http.StatusConflict, "stock product is less than quantity")
	if got != nil {
		t.Fatalf("cart = %#v, want nil", got)
	}
}

func TestCartServiceAddItemReturnsNotFoundWhenProductMissing(t *testing.T) {
	cart := &cartModel.Cart{ID: "cart-1", CustomerID: "customer-1"}
	service := NewCartService(
		&mockCartRepository{
			getOrCreateCartFn: func(context.Context, string) (*cartModel.Cart, error) {
				return cart, nil
			},
		},
		&mockCartItemRepository{
			getCartItemsFn: func(context.Context, string, int) ([]*cartModel.CartItem, error) {
				return []*cartModel.CartItem{}, nil
			},
		},
		&mockProductService{
			findByIDFn: func(context.Context, string) (*catalogDto.ProductRes, error) {
				return nil, nil
			},
		},
		&mockInventoryService{
			findByProductIDFn: func(context.Context, string) (*catalogModel.Inventory, error) {
				return &catalogModel.Inventory{StockQuantity: 10}, nil
			},
		},
	)

	got, err := service.AddItem(context.Background(), "customer-1", "product-1", 1)
	assertCartAppError(t, err, http.StatusNotFound, "product not found")
	if got != nil {
		t.Fatalf("cart = %#v, want nil", got)
	}
}

func TestCartServiceUpdateItem(t *testing.T) {
	cart := &cartModel.Cart{ID: "cart-1", CustomerID: "customer-1"}
	cartRepo := &mockCartRepository{
		getCartByCustomerIDFn: func(context.Context, string) (*cartModel.Cart, error) {
			return cart, nil
		},
	}
	itemRepo := &mockCartItemRepository{
		getCartItemsFn: func(_ context.Context, _ string, call int) ([]*cartModel.CartItem, error) {
			if call == 1 {
				return []*cartModel.CartItem{{ProductID: "product-1", Quantity: 2, PriceSnapshot: 11000}}, nil
			}
			return []*cartModel.CartItem{{ProductID: "product-1", Quantity: 4, PriceSnapshot: 11000}}, nil
		},
		upsertItemFn: func(_ context.Context, _ string, productID string, quantity int, priceSnapshot float64) error {
			if productID != "product-1" || quantity != 4 || priceSnapshot != 11000 {
				t.Fatalf("upsert args = %q %d %.2f, want product-1 4 11000", productID, quantity, priceSnapshot)
			}
			return nil
		},
	}
	inventorySvc := &mockInventoryService{
		findByProductIDFn: func(context.Context, string) (*catalogModel.Inventory, error) {
			return &catalogModel.Inventory{StockQuantity: 5}, nil
		},
	}

	got, err := NewCartService(cartRepo, itemRepo, &mockProductService{}, inventorySvc).
		UpdateItem(context.Background(), "customer-1", "product-1", 4)
	if err != nil {
		t.Fatalf("UpdateItem returned error: %v", err)
	}
	assertCartResponse(t, got, "cart-1", "customer-1", 44000, []dto.CartItemRes{
		{ProductID: "product-1", Quantity: 4, PriceSnapshot: 11000, Subtotal: 44000},
	})
}

func TestCartServiceUpdateItemReturnsNotFoundWhenCartMissing(t *testing.T) {
	service := NewCartService(
		&mockCartRepository{
			getCartByCustomerIDFn: func(context.Context, string) (*cartModel.Cart, error) {
				return nil, nil
			},
		},
		&mockCartItemRepository{},
		&mockProductService{},
		&mockInventoryService{},
	)

	got, err := service.UpdateItem(context.Background(), "customer-1", "product-1", 1)
	assertCartAppError(t, err, http.StatusNotFound, "cart not found")
	if got != nil {
		t.Fatalf("cart = %#v, want nil", got)
	}
}

func TestCartServiceUpdateItemReturnsNotFoundWhenItemMissing(t *testing.T) {
	service := NewCartService(
		&mockCartRepository{
			getCartByCustomerIDFn: func(context.Context, string) (*cartModel.Cart, error) {
				return &cartModel.Cart{ID: "cart-1", CustomerID: "customer-1"}, nil
			},
		},
		&mockCartItemRepository{
			getCartItemsFn: func(context.Context, string, int) ([]*cartModel.CartItem, error) {
				return []*cartModel.CartItem{}, nil
			},
		},
		&mockProductService{},
		&mockInventoryService{},
	)

	got, err := service.UpdateItem(context.Background(), "customer-1", "product-1", 1)
	assertCartAppError(t, err, http.StatusNotFound, "cart item not found")
	if got != nil {
		t.Fatalf("cart = %#v, want nil", got)
	}
}

func TestCartServiceDeleteItem(t *testing.T) {
	cart := &cartModel.Cart{ID: "cart-1", CustomerID: "customer-1"}
	cartRepo := &mockCartRepository{
		getCartByCustomerIDFn: func(context.Context, string) (*cartModel.Cart, error) {
			return cart, nil
		},
	}
	itemRepo := &mockCartItemRepository{
		deleteItemFn: func(_ context.Context, cartID string, productID string) error {
			if cartID != cart.ID || productID != "product-1" {
				t.Fatalf("delete args = %q %q, want cart-1 product-1", cartID, productID)
			}
			return nil
		},
	}

	if err := NewCartService(cartRepo, itemRepo, &mockProductService{}, &mockInventoryService{}).
		DeleteItem(context.Background(), "customer-1", "product-1"); err != nil {
		t.Fatalf("DeleteItem returned error: %v", err)
	}
}

func TestCartServiceDeleteItemReturnsNotFoundWhenCartMissing(t *testing.T) {
	service := NewCartService(
		&mockCartRepository{
			getCartByCustomerIDFn: func(context.Context, string) (*cartModel.Cart, error) {
				return nil, nil
			},
		},
		&mockCartItemRepository{},
		&mockProductService{},
		&mockInventoryService{},
	)

	err := service.DeleteItem(context.Background(), "customer-1", "product-1")
	assertCartAppError(t, err, http.StatusNotFound, "cart not found")
}

func TestCartServiceClearItems(t *testing.T) {
	cart := &cartModel.Cart{ID: "cart-1", CustomerID: "customer-1"}
	cartRepo := &mockCartRepository{
		getCartByCustomerIDFn: func(context.Context, string) (*cartModel.Cart, error) {
			return cart, nil
		},
	}
	itemRepo := &mockCartItemRepository{
		clearItemsFn: func(_ context.Context, cartID string) error {
			if cartID != cart.ID {
				t.Fatalf("cart id = %q, want %q", cartID, cart.ID)
			}
			return nil
		},
	}

	if err := NewCartService(cartRepo, itemRepo, &mockProductService{}, &mockInventoryService{}).
		ClearItems(context.Background(), "customer-1"); err != nil {
		t.Fatalf("ClearItems returned error: %v", err)
	}
}

func TestCartServiceClearItemsRejectsInvalidCustomerID(t *testing.T) {
	service := NewCartService(&mockCartRepository{}, &mockCartItemRepository{}, &mockProductService{}, &mockInventoryService{})

	err := service.ClearItems(context.Background(), "")
	assertCartAppError(t, err, http.StatusBadRequest, "invalid parameter customer id")
}

func assertCartAppError(t *testing.T, err error, code int, message string) {
	t.Helper()

	var appErr *response.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %T %v, want *response.AppError", err, err)
	}
	if appErr.Code != code {
		t.Fatalf("code = %d, want %d", appErr.Code, code)
	}
	if appErr.Message != message {
		t.Fatalf("message = %q, want %q", appErr.Message, message)
	}
}

func assertCartResponse(t *testing.T, got *dto.CartRes, cartID string, customerID string, total float64, items []dto.CartItemRes) {
	t.Helper()

	if got == nil {
		t.Fatal("cart response is nil")
	}
	if got.CartID != cartID {
		t.Fatalf("cart id = %q, want %q", got.CartID, cartID)
	}
	if got.CustomerID != customerID {
		t.Fatalf("customer id = %q, want %q", got.CustomerID, customerID)
	}
	if got.Total != total {
		t.Fatalf("total = %.2f, want %.2f", got.Total, total)
	}
	if len(got.Items) != len(items) {
		t.Fatalf("items count = %d, want %d", len(got.Items), len(items))
	}
	for i, want := range items {
		if got.Items[i] != want {
			t.Fatalf("item[%d] = %#v, want %#v", i, got.Items[i], want)
		}
	}
}
