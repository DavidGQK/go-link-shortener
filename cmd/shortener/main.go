package main

import (
	"github.com/DavidGQK/go-link-shortener/internal/app"
	"log"
)

func main() {
	if err := app.RunServer(); err != nil {
		log.Fatal(err)
	}
}
