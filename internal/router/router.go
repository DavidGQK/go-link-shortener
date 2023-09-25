package router

import (
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/DavidGQK/go-link-shortener/internal/middleware"
	"github.com/DavidGQK/go-link-shortener/internal/server"
	"github.com/go-chi/chi/v5"
)

func NewRouter(s server.Server) chi.Router {
	r := chi.NewRouter()
	r.Use(logger.LoggingMiddleware, middleware.GzipMiddleware)
	r.Get("/{id}", s.GetContent)
	r.Post("/", s.PostShortenLink)
	r.Post("/api/shorten", s.PostAPIShortenLink)

	return r
}
