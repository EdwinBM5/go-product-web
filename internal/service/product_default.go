package service

import "github.com/edwinbm5/go-product-web/internal"

type ProductDefault struct {
	repository internal.ProductRepository
}

func NewDefaultProduct(repository internal.ProductRepository) *ProductDefault {
	return &ProductDefault{
		repository: repository,
	}
}

func (p *ProductDefault) GetAll() (products []internal.Product, err error) {
	products, err = p.repository.GetAll()
	return
}

func (p *ProductDefault) GetByID(id int) (product internal.Product, err error) {
	product, err = p.repository.GetByID(id)
	return
}

func (p *ProductDefault) Create(product *internal.Product) (err error) {
	err = p.repository.Create(product)
	return
}
