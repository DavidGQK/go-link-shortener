package main

import (
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/DavidGQK/go-link-shortener/internal/router"
	"github.com/DavidGQK/go-link-shortener/internal/server"
	"github.com/DavidGQK/go-link-shortener/internal/storage"
	"net/http"
)

func runServer(cfg *config.Config) error {
	db, err := storage.New(cfg.Filename)
	if err != nil {
		return err
	}

	s := server.New(cfg, db)
	if err := logger.Initialize(cfg.LoggingLevel); err != nil {
		return err
	}

	r := router.NewRouter(s)

	logger.Log.Infow("server start", "address", cfg.ServerURL)
	return http.ListenAndServe(cfg.ServerURL, r)
}

func main() {
	cfg := config.GetConfig()
	if err := runServer(cfg); err != nil {
		panic(err)
	}
}
