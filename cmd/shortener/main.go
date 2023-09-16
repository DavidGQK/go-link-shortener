package main

import (
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/DavidGQK/go-link-shortener/internal/router"
	"github.com/DavidGQK/go-link-shortener/internal/server"
	"github.com/DavidGQK/go-link-shortener/internal/storage"
	"log"
	"net/http"
)

func RunServer(cfg *config.Config, serverStorage *storage.Storage) error {
	s := server.NewServer(cfg.ShortURLBase, serverStorage)
	r := router.NewRouter(s)
	return http.ListenAndServe(cfg.ServerURL, r)
}

func main() {
	db := storage.New()
	cfg := config.GetConfig()
	if err := RunServer(cfg, db); err != nil {
		log.Fatal(err)
	}
}
