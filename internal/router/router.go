package router

import (
	"github.com/DavidGQK/go-link-shortener/internal/server"
	"github.com/go-chi/chi/v5"
)

func NewRouter(s server.Server) chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", s.GetContent)
	r.Post("/", s.PostShortLink)

	return r
}
