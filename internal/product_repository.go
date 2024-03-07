package internal

type ProductRepository interface {
	Save(product *Product) (err error)
}
