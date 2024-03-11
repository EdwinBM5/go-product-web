package repository

import (
	"fmt"

	"github.com/edwinbm5/go-product-web/internal"
	"github.com/edwinbm5/go-product-web/internal/platform/tools"
)

// ProductSlice is a repository that stores products in a slice
type ProductSlice struct {
	db     []internal.Product
	lastID int
}

// NewProductSlice creates a new ProductSlice
func NewProductSlice(db []internal.Product, lastID int) *ProductSlice {
	if db == nil {
		db = make([]internal.Product, 0)
	}

	return &ProductSlice{
		db:     db,
		lastID: lastID,
	}
}

// GetAll returns all the products in the database
func (p *ProductSlice) GetAll() (products []internal.Product, err error) {
	if len(p.db) == 0 {
		err = internal.ErrProductsEmpty
		return
	}

	products = p.db

	return
}

// GetByID returns a product by its ID
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

// Creates a new product in the database
func (p *ProductSlice) Create(product *internal.Product) (err error) {
	for _, pr := range p.db {
		if pr.CodeValue == (*product).CodeValue {
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

// Updates a product in the database or creates it if it does not exist
func (p *ProductSlice) UpdateAndCreate(product *internal.Product) (err error) {
	for _, pr := range p.db {
		if pr.ID == product.ID {
			err = p.Update(product.ID, map[string]any{
				"CodeValue":   product.CodeValue,
				"Name":        product.Name,
				"Price":       product.Price,
				"Expiration":  product.Expiration,
				"IsPublished": product.IsPublished,
				"Quantity":    product.Quantity,
				"ID":          product.ID,
			})

			return
		}
	}

	err = p.Create(product)

	return
}

// Updates a product in the database
func (p *ProductSlice) Update(id int, fields map[string]any) (err error) {
	productExist := false
	productIndex := 0

	for index, pr := range p.db {
		if pr.ID == id {
			productExist = true
			productIndex = index
		}
	}

	if !productExist {
		err = internal.ErrProductNotFound
		err = fmt.Errorf("%w: The product with ID %d does not exist", err, id)
		return
	}

	product := p.db[productIndex]

	err = tools.ParseDate(product.Expiration)
	if err != nil {
		err = tools.ErrInvalidDate
		return
	}

	for key := range fields {
		switch key {
		case "CodeValue", "codevalue":
			for _, pr := range p.db {
				if pr.CodeValue == fields["codevalue"] && pr.ID != id {
					err = internal.ErrProductDuplicated
					err = fmt.Errorf("%w: The Code value %s already exists", err, product.CodeValue)
					return
				}
			}
			product.CodeValue = fields[key].(string)
		case "Name", "name":
			product.Name = fields[key].(string)
		case "Expiration", "expiration":
			product.Expiration = fields[key].(string)
		case "Price", "price":
			product.Price = fields[key].(float64)
		case "Quantity", "quantity":
			product.Quantity = fields[key].(int)
		case "IsPublished", "ispublished":
			product.IsPublished = fields[key].(bool)
		default:
		}
	}

	p.db[productIndex] = product

	return
}

// Deletes a product from the database
func (p *ProductSlice) Delete(id int) (err error) {
	for index, product := range p.db {
		if product.ID == id {
			p.db = append(p.db[:index], p.db[index+1:]...)
			return
		}
	}

	err = internal.ErrProductNotFound

	return
}
