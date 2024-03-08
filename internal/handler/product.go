package handler

import (
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/edwinbm5/go-product-web/internal"
)

type DefaultProduct struct {
	sv internal.ProductService
}

func NewDefaultProduct(sv internal.ProductService) *DefaultProduct {
	return &DefaultProduct{
		sv: sv,
	}
}

type ProductJSON struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type ProductRequestBody struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func (d *DefaultProduct) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := d.sv.GetAll()
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
			return
		}

		response.JSON(w, http.StatusOK, map[string]interface{}{
			"message":  "total products: " + strconv.Itoa(len(products)),
			"products": products,
		})
	}
}

func (d *DefaultProduct) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (d *DefaultProduct) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
