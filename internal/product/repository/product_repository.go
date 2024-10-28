package repository

import (
	model "Dzaakk/simple-commerce/internal/product/models"
	modelCart "Dzaakk/simple-commerce/internal/shopping_cart/models"
)

type ProductRepository interface {
	Create(data model.TProduct) (*model.TProduct, error)
	Update(data model.TProduct) error
	FindByCategoryId(categoryId int) ([]*model.TProduct, error)
	FindById(id int) (*model.TProduct, error)
	FindByName(name string) (*model.TProduct, error)
	GetPriceById(id int) (*float32, error)
	GetStockById(id int) (int, error)
	UpdateStock(listData []*modelCart.TCartItemDetail, name string) error
}
