package internal

import "errors"

type Product struct {
	ID          int
	Name        string
	Quantity    int
	CodeValue   string
	IsPublished bool
	Expiration  string
	Price       float64
}

var (
	ErrProductNotFound     = errors.New("Product not found")
	ErrProductDuplicated   = errors.New("Product already exists")
	ErrProductInternal     = errors.New("Product can't be processed")
	ErrProductsEmpty       = errors.New("Products are empty")
	ErrProductInvalidField = errors.New("Product field is invalid")
)
