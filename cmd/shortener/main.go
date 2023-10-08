package main

import (
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/DavidGQK/go-link-shortener/internal/router"
	"github.com/DavidGQK/go-link-shortener/internal/server"
	"github.com/DavidGQK/go-link-shortener/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func runServer(cfg *config.Config) error {
	var dataWr *storage.DataWriter
	var err error

	if cfg.Filename != "" {
		file, err := os.OpenFile(cfg.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.Log.Error("open file error", zap.Error(err))
			return err
		}

		dataWr, err = storage.NewDataWriter(file)
		if err != nil {
			logger.Log.Error("creating a new data writer error", zap.Error(err))
			return err
		}
		defer dataWr.Close()
	}

	st, err := storage.New(cfg.Filename, cfg.DBConnData, dataWr)
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
