package router

import (
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/DavidGQK/go-link-shortener/internal/middleware"
	"github.com/DavidGQK/go-link-shortener/internal/server"
	"github.com/go-chi/chi/v5"
)

func NewRouter(s server.Server) chi.Router {
	r := chi.NewRouter()
	r.Use(logger.Middleware, middleware.GzipMiddleware)
	r.Get("/{id}", s.GetContent)
	r.Post("/", s.PostShortenLink)
	r.Post("/api/shorten", s.PostAPIShortenLink)
	r.Get("/ping", s.Ping)
	r.Post("/api/shorten/batch", s.PostAPIShortenBatch)

	return r
}
