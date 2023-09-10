package main

import (
	"github.com/DavidGQK/go-link-shortener/internal/router"
	"github.com/DavidGQK/go-link-shortener/internal/server"
	"github.com/DavidGQK/go-link-shortener/internal/storage"
	"log"
	"net/http"
)

func RunServer() error {
	serverStorage := storage.NewStorage()
	s := server.NewServer(serverStorage)
	r := router.NewRouter(s)
	return http.ListenAndServe(s.Config.ServerURL, r)
}

func main() {
	if err := RunServer(); err != nil {
		log.Fatal(err)
	}
}
