package repositories

import (
	model "Dzaakk/simple-commerce/internal/product/models"
	modelCart "Dzaakk/simple-commerce/internal/shopping_cart/models"
	"context"
	"database/sql"
)

type ProductRepository interface {
	Create(ctx context.Context, data model.TProduct) (*model.TProduct, error)
	Update(ctx context.Context, data model.TProduct) (int64, error)
	FindByProductID(ctx context.Context, productID int) (*model.TProduct, error)
	FindByProductName(ctx context.Context, productName string) (*model.TProduct, error)
	SetStockByProductID(ctx context.Context, productID int, stock int) (int64, error)
	GetPriceByProductID(ictx context.Context, productID int) (float32, error)
	GetStockByProductID(ctx context.Context, productID int) (int, error)
	UpdateStock(ctx context.Context, listData []*modelCart.TCartItemDetail, name string) error
	UpdateStockWithTx(ctx context.Context, tx *sql.Tx, listData []*modelCart.TCartItemDetail) ([]*int, error)
	FindProductByFilters(ctx context.Context, sellerID, categoryID *int) ([]*model.TProduct, error)
}
