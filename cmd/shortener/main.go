package main

import (
	"github.com/DavidGQK/go-link-shortener/internal/router"
	"github.com/DavidGQK/go-link-shortener/internal/server"
	"github.com/DavidGQK/go-link-shortener/internal/storage"
	"log"
	"net/http"
)

func RunServer(serverStorage *storage.Storage) error {
	s := server.NewServer(serverStorage)
	r := router.NewRouter(s)
	return http.ListenAndServe(s.Config.ServerURL, r)
}

func main() {
	db := storage.New()
	if err := RunServer(db); err != nil {
		log.Fatal(err)
	}
}
