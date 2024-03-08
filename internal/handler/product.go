package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/edwinbm5/go-product-web/internal"
	"github.com/edwinbm5/go-product-web/internal/platform/tools"
	"github.com/go-chi/chi/v5"
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
			response.JSON(w, http.StatusInternalServerError, map[string]any{"message": err.Error()})
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message":  "Total products: " + strconv.Itoa(len(products)),
			"products": products,
		})
	}
}

func (d *DefaultProduct) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{"message": "Invalid ID"})
			return
		}

		product, err := d.sv.GetByID(idInt)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, map[string]any{"message": err.Error()})
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "Product found",
			"product": product,
		})
	}
}

func (d *DefaultProduct) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request - read into bytes
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid request body",
			})

			return
		}

		// parse to map (dynamic)
		bodyMap := map[string]any{}
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid request body",
			})

			return
		}

		// validate required fields
		if err := tools.CheckFieldExistance(bodyMap, "name", "quantity", "code_value", "expiration", "price"); err != nil {
			var fieldError *tools.FieldError
			if errors.As(err, &fieldError) {
				response.JSON(w, http.StatusBadRequest, map[string]any{
					"message": fmt.Sprintf("%s is required", fieldError.Field),
				})

				return
			}
		}

		// parse json to struct (static)
		var body ProductRequestBody
		if err := json.Unmarshal(bytes, &body); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid request body",
			})

			return
		}

		// validate the product
		product := internal.Product{
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}

		// create the product
		if err := d.sv.Create(&product); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductDuplicated):
				response.JSON(w, http.StatusConflict, map[string]any{
					"message": "Product already exists",
				})
			case errors.Is(err, tools.ErrInvalidDay):
				response.JSON(w, http.StatusBadRequest, map[string]any{
					"message": "Invalid day on date",
				})
			case errors.Is(err, tools.ErrInvalidMonth):
				response.JSON(w, http.StatusBadRequest, map[string]any{
					"message": "Invalid month on date",
				})
			case errors.Is(err, tools.ErrInvalidYear):
				response.JSON(w, http.StatusBadRequest, map[string]any{
					"message": "Invalid year on date",
				})
			default:
				response.JSON(w, http.StatusInternalServerError, map[string]any{
					"message": "Internal server error",
				})
			}

			return
		}

		// response
		data := ProductJSON{
			ID:          product.ID,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "Product created successfully",
			"data":    data,
		})
	}
}
