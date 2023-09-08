package app

import (
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/DavidGQK/go-link-shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RunServer() error {
	config.ParseFlags()

	r := chi.NewRouter()
	r.Get("/{id}", handlers.ProcessGET)
	r.Post("/", handlers.ProcessPOST)

	return http.ListenAndServe(config.ServerURL, r)
}
