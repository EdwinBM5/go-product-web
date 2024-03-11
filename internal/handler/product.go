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
	"github.com/edwinbm5/go-product-web/internal/auth"
	"github.com/edwinbm5/go-product-web/internal/platform/tools"
	"github.com/go-chi/chi/v5"
)

type DefaultProduct struct {
	sv internal.ProductService
	au auth.Auth
}

func NewDefaultProduct(sv internal.ProductService, au auth.Auth) *DefaultProduct {
	return &DefaultProduct{
		sv: sv,
		au: au,
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

// GetAll is a handler for get all the products in the database
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

// GetByID is a handler for get by ID a product
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

// Create is a handler for Create a new product in the database
func (d *DefaultProduct) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Validate authentication
		token := r.Header.Get("token")
		if token != d.au.GetToken() {
			response.JSON(w, http.StatusUnauthorized, map[string]any{"message": "Unauthorized"})
			return
		}

		// Request - Read into bytes
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid request body",
			})

			return
		}

		// Parse to map (dynamic)
		bodyMap := map[string]any{}
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid request body",
			})

			return
		}

		// Validate required fields
		if err := tools.CheckFieldExistance(bodyMap, "name", "quantity", "code_value", "expiration", "price"); err != nil {
			var fieldError *tools.FieldError
			if errors.As(err, &fieldError) {
				response.JSON(w, http.StatusBadRequest, map[string]any{
					"message": fmt.Sprintf("%s is required", fieldError.Field),
				})

				return
			}
		}

		// Parse json to struct (static)
		var body ProductRequestBody
		if err := json.Unmarshal(bytes, &body); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid request body",
			})

			return
		}

		// Validate the product
		product := internal.Product{
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}

		// Create the product
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

		// Response
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

// UpdateAndCreate is a handler for update or create a product in the database if not exists
func (d *DefaultProduct) UpdateAndCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate authentication
		token := r.Header.Get("token")
		if token != d.au.GetToken() {
			response.JSON(w, http.StatusUnauthorized, map[string]any{"message": "Unauthorized"})
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid ID",
			})

			return
		}

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid request body",
			})

			return
		}

		bodyMap := map[string]any{}
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid request body",
			})
		}

		if err := tools.CheckFieldExistance(bodyMap, "name", "quantity", "code_value", "expiration", "price"); err != nil {
			var fieldError *tools.FieldError
			if errors.As(err, &fieldError) {
				response.JSON(w, http.StatusBadRequest, map[string]any{
					"message": fmt.Sprintf("%s is required", fieldError.Field),
				})

				return
			}
		}

		var body ProductRequestBody
		if err := json.Unmarshal(bytes, &body); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid request body",
			})

			return
		}

		product := internal.Product{
			ID:          id,
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}

		if err := d.sv.UpdateAndCreate(&product); err != nil {
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

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "Product updated successfully",
			"data": ProductJSON{
				ID:          product.ID,
				Name:        product.Name,
				Quantity:    product.Quantity,
				CodeValue:   product.CodeValue,
				IsPublished: product.IsPublished,
				Expiration:  product.Expiration,
				Price:       product.Price,
			},
		})
	}
}

// Update is a handler for update a product in the database
func (d *DefaultProduct) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate authentication
		token := r.Header.Get("token")
		if token != d.au.GetToken() {
			response.JSON(w, http.StatusUnauthorized, map[string]any{"message": "Unauthorized"})
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid ID",
			})

			return
		}

		bodyMap := make(map[string]any)
		if err := json.NewDecoder(r.Body).Decode(&bodyMap); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid request body",
			})

			return
		}

		if name, ok := bodyMap["name"]; ok {
			_, ok := name.(string)
			if !ok {
				response.JSON(w, http.StatusBadRequest, map[string]any{
					"message": "Invalid name",
				})

				return
			}
		}

		if quantity, ok := bodyMap["quantity"]; ok {
			_, ok := quantity.(int)
			if !ok {
				response.JSON(w, http.StatusBadRequest, map[string]any{
					"message": "Invalid quantity",
				})

				return
			}
		}

		if codeValue, ok := bodyMap["code_value"]; ok {
			_, ok := codeValue.(string)
			if !ok {
				response.JSON(w, http.StatusBadRequest, map[string]any{
					"message": "Invalid code value",
				})

				return
			}
		}

		if isPublished, ok := bodyMap["is_published"]; ok {
			_, ok := isPublished.(bool)
			if !ok {
				response.JSON(w, http.StatusBadRequest, map[string]any{
					"message": "Invalid is published",
				})

				return
			}
		}

		if expiration, ok := bodyMap["expiration"]; ok {
			_, ok := expiration.(string)
			if !ok {
				response.JSON(w, http.StatusBadRequest, map[string]any{
					"message": "Invalid expiration",
				})

				return
			}
		}

		if price, ok := bodyMap["price"]; ok {
			_, ok := price.(float64)
			if !ok {
				response.JSON(w, http.StatusBadRequest, map[string]any{
					"message": "Invalid price",
				})

				return
			}
		}

		if err := d.sv.Update(id, bodyMap); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.JSON(w, http.StatusNotFound, map[string]any{
					"message": "Product not found",
				})
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

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "Product updated successfully",
		})
	}
}

// Delete is a handler for delete a product in the database
func (d *DefaultProduct) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid ID",
			})

			return
		}

		if err := d.sv.Delete(id); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.JSON(w, http.StatusNotFound, map[string]any{
					"message": "Product not found",
				})
			default:
				response.JSON(w, http.StatusInternalServerError, map[string]any{
					"message": "Internal server error",
				})

			}

			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "Product deleted successfully",
		})

	}
}
