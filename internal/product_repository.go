package internal

type ProductRepository interface {
	GetAll() (products []Product, err error)
	GetByID(id int) (product Product, err error)
	Create(product *Product) (err error)
}
