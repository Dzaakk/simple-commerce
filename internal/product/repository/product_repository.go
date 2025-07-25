package repository

import (
	"Dzaakk/simple-commerce/internal/product/model"
	cartModel "Dzaakk/simple-commerce/internal/shopping_cart/model"
	"context"
	"database/sql"
)

type ProductRepository interface {
	Create(ctx context.Context, data model.TProduct) (*model.TProduct, error)
	Update(ctx context.Context, data model.TProduct) (int64, error)
	FindByID(ctx context.Context, id int) (*model.TProduct, error)
	FindByProductName(ctx context.Context, productName string) (*model.TProduct, error)
	UpdateStock(ctx context.Context, listData []*cartModel.TCartItemDetail, name string) error
	UpdateStockWithTx(ctx context.Context, tx *sql.Tx, listData []*cartModel.TCartItemDetail) ([]*int, error)
	FindByFilters(ctx context.Context, params model.ProductFilter) ([]*model.TProduct, error)
	// SetStockByID(ctx context.Context, id int, stock int) (int64, error)
	// GetPriceByID(ictx context.Context, id int) (float32, error)
	// GetStockByID(ctx context.Context, id int) (int, error)
	// FindProductByFilters(ctx context.Context, sellerID, categoryID *int) ([]*model.TProduct, error)
}
