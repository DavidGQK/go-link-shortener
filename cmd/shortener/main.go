package main

import (
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/DavidGQK/go-link-shortener/internal/router"
	"github.com/DavidGQK/go-link-shortener/internal/server"
	"github.com/DavidGQK/go-link-shortener/internal/storage/initstorage"
	"go.uber.org/zap"
	"net/http"
)

func runServer(cfg *config.Config) error {
	st, err := initstorage.NewStorage(cfg.Filename, cfg.DBConnData)
	if err != nil {
		return err
	}

	if cfg.Filename != "" {
		if err := st.Restore(); err != nil {
			logger.Log.Error("restore storage error", zap.Error(err))
		}
	}

	s := server.New(cfg, st)
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
