package config

import (
	"flag"
	"os"
)

var (
	ServerURL    string
	ShortURLBase string
)

func ParseFlags() {
	flag.StringVar(&ServerURL, "a", "localhost:8080", "url where server runs on")
	flag.StringVar(&ShortURLBase, "b", "http://localhost:8080", "base url for shortened link")
	flag.Parse()

	if envServerURL := os.Getenv("SERVER_ADDRESS"); envServerURL != "" {
		ServerURL = envServerURL
	}

	if envServerURLBase := os.Getenv("BASE_URL"); envServerURLBase != "" {
		ShortURLBase = envServerURLBase
	}
}
