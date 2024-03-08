package main

import (
	"fmt"
	"net/http"

	"github.com/edwinbm5/go-product-web/internal/handler"
	"github.com/edwinbm5/go-product-web/internal/repository"

	"github.com/edwinbm5/go-product-web/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	repository := repository.NewProductSlice(nil, 0)
	service := service.NewDefaultProduct(repository)
	handler := handler.NewDefaultProduct(service)

	router := chi.NewRouter()
	router.Route("/products", func(r chi.Router) {
		r.Get("/", handler.GetAll())
		r.Post("/", handler.Create())
		r.Get("/{id}", handler.GetByID())
	})

	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		fmt.Println(err)
		return
	}
}
