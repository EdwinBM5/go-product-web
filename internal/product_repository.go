package internal

type ProductRepository interface {
	GetAll() (products []Product, err error)
	GetByID(id int) (product Product, err error)
	Create(product *Product) (err error)
	UpdateAndCreate(product *Product) (err error)
	Update(id int, fields map[string]any) (err error)
	Delete(id int) (err error)
}
