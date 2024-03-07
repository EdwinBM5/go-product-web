package internal

type ProductService interface {
	Save(product *Product) (err error)
}
