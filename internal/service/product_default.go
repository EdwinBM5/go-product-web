package service

import "github.com/edwinbm5/go-product-web/internal"

type ProductDefault struct {
	repository internal.ProductRepository
}

// NewDefaultProduct creates a new ProductDefault service
func NewDefaultProduct(repository internal.ProductRepository) *ProductDefault {
	return &ProductDefault{
		repository: repository,
	}
}

// GetAll returns all the products in the database
func (p *ProductDefault) GetAll() (products []internal.Product, err error) {
	products, err = p.repository.GetAll()
	return
}

// GetByID returns a product by its ID
func (p *ProductDefault) GetByID(id int) (product internal.Product, err error) {
	product, err = p.repository.GetByID(id)
	return
}

// Creates a new product in the database
func (p *ProductDefault) Create(product *internal.Product) (err error) {
	err = p.repository.Create(product)
	return
}

// Updates a product in the database, if not exists, creates it
func (p *ProductDefault) UpdateAndCreate(product *internal.Product) (err error) {
	err = p.repository.UpdateAndCreate(product)
	return
}

// Updates a product in the database
func (p *ProductDefault) Update(id int, fields map[string]any) (err error) {
	err = p.repository.Update(id, fields)
	return
}

func (p *ProductDefault) Delete(id int) (err error) {
	err = p.repository.Delete(id)
	return
}
