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
	r.Get("/{id}", s.CookieMiddleware(s.GetContent))
	r.Post("/", s.CookieMiddleware(s.PostShortenLink))
	r.Post("/api/shorten", s.CookieMiddleware(s.PostAPIShortenLink))
	r.Get("/ping", s.CookieMiddleware(s.Ping))
	r.Post("/api/shorten/batch", s.CookieMiddleware(s.PostAPIShortenBatch))
	r.Get("/api/user/urls", s.CookieMiddleware(s.GetUserUrlsAPI))
	r.Delete("/api/user/urls", s.CookieMiddleware(s.DeleteUserUrls))

	return r
}
