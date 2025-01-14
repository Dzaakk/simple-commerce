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
	FindByCategoryId(ctx context.Context, categoryId int) ([]*model.TProduct, error)
	FindById(ctx context.Context, id int) (*model.TProduct, error)
	FindByName(ctx context.Context, name string) (*model.TProduct, error)
	SetStockById(ctx context.Context, id int, stock int) (int64, error)
	GetPriceById(ictx context.Context, d int) (float32, error)
	GetStockById(ctx context.Context, id int) (int, error)
	UpdateStock(ctx context.Context, listData []*modelCart.TCartItemDetail, name string) error
	UpdateStockWithTx(ctx context.Context, tx *sql.Tx, listData []*modelCart.TCartItemDetail) ([]*int, error)

	//Seller
	FindBySellerId(ctx context.Context, sellerId int) ([]*model.TProduct, error)
	FindBySellerIdAndCategoryId(ctx context.Context, sellerId, categoryId int) ([]*model.TProduct, error)
}
