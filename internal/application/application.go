package application

import (
	"fmt"
	"net/http"

	"github.com/edwinbm5/go-product-web/internal/auth"
	"github.com/edwinbm5/go-product-web/internal/handler"
	"github.com/edwinbm5/go-product-web/internal/repository"
	"github.com/edwinbm5/go-product-web/internal/service"
	"github.com/go-chi/chi/v5"
)

type DefaultApp struct {
	Title    string
	Color    string
	FilePath string
	Token    string
}

type ConfigDefaultApp struct {
	Title    string `json:"title"`
	Color    string `json:"color"`
	FilePath string `json:"file_path"`
	Token    string `json:"token"`
}

func NewDefaultApp(cfg ConfigDefaultApp) *DefaultApp {
	fmt.Println(cfg.Color, cfg.Title, cfg.FilePath)
	if cfg.Title == "" {
		cfg.Title = "Generic App"
	}

	return &DefaultApp{
		Title:    cfg.Title,
		Color:    cfg.Color,
		FilePath: cfg.FilePath,
		Token:    cfg.Token,
	}
}

func (d *DefaultApp) Run() {
	fmt.Println("Running application...")
	fmt.Printf("Title: %s\n", d.Title)

	auth := auth.NewAuthDefault(d.Token)

	repository := repository.NewProductSlice(nil, 0)
	service := service.NewDefaultProduct(repository)
	handler := handler.NewDefaultProduct(service, auth)

	router := chi.NewRouter()
	router.Route("/products", func(r chi.Router) {
		r.Get("/", handler.GetAll())
		r.Post("/", handler.Create())
		r.Get("/{id}", handler.GetByID())
		r.Patch("/{id}", handler.Update())
		r.Put("/{id}", handler.UpdateAndCreate())
		r.Delete("/{id}", handler.Delete())
	})

	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		fmt.Println(err)
		return
	}
}
