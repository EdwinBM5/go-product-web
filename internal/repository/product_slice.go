package repository

import (
	"fmt"

	"github.com/edwinbm5/go-product-web/internal"
	"github.com/edwinbm5/go-product-web/internal/platform/tools"
)

type ProductSlice struct {
	db     []internal.Product
	lastID int
}

func NewProductSlice(db []internal.Product, lastID int) *ProductSlice {
	if db == nil {
		db = make([]internal.Product, 0)
	}

	return &ProductSlice{
		db:     db,
		lastID: lastID,
	}
}

func (p *ProductSlice) GetAll() (products []internal.Product, err error) {
	if len(p.db) == 0 {
		err = internal.ErrProductsEmpty
		return
	}

	products = p.db

	return
}

func (p *ProductSlice) GetByID(id int) (product internal.Product, err error) {
	for _, prod := range p.db {
		if prod.ID == id {
			product = prod
			return
		}
	}

	err = internal.ErrProductNotFound
	err = fmt.Errorf("%w: The product with ID %d does not exist", err, id)
	return
}

func (p *ProductSlice) Create(product *internal.Product) (err error) {
	for _, p := range (*p).db {
		if p.CodeValue == (*product).CodeValue {
			err = internal.ErrProductDuplicated
			err = fmt.Errorf("%w: The Code value %s already exists", err, (*product).CodeValue)
			return
		}
	}

	err = tools.ParseDate(product.Expiration)
	if err != nil {
		err = fmt.Errorf("%w: The expiration date %s is not valid", err, product.Expiration)
		return
	}

	p.lastID++
	product.ID = p.lastID

	p.db = append(p.db, *product)

	return
}
