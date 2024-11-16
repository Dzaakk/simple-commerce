package repository

import (
	model "Dzaakk/simple-commerce/internal/product/models"
	modelCart "Dzaakk/simple-commerce/internal/shopping_cart/models"
	"database/sql"
)

type ProductRepository interface {
	Create(data model.TProduct) (*model.TProduct, error)
	Update(data model.TProduct) error
	FindByCategoryId(categoryId int) ([]*model.TProduct, error)
	FindById(id int) (*model.TProduct, error)
	FindByName(name string) (*model.TProduct, error)
	SetStockById(id int, stock int) error
	GetPriceById(id int) (*float32, error)
	GetStockById(id int) (int, error)
	UpdateStock(listData []*modelCart.TCartItemDetail, name string) error
	UpdateStockWithTx(tx *sql.Tx, listData []*modelCart.TCartItemDetail) ([]*int, error)
}