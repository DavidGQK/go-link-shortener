package main

import (
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/DavidGQK/go-link-shortener/internal/router"
	"github.com/DavidGQK/go-link-shortener/internal/server"
	"github.com/DavidGQK/go-link-shortener/internal/storage"
	"go.uber.org/zap"
	"net/http"
)

func runServer(cfg *config.Config, serverStorage *storage.Storage) error {
	s := server.NewServer(cfg.ShortURLBase, serverStorage)
	if err := logger.Initialize(cfg.LoggingLevel); err != nil {
		return err
	}

	r := router.NewRouter(s)

	logger.Log.Info("server start", zap.String("address", cfg.ServerURL))
	return http.ListenAndServe(cfg.ServerURL, r)
}

func main() {
	db := storage.New()
	cfg := config.GetConfig()
	if err := runServer(cfg, db); err != nil {
		panic(err)
	}
}
