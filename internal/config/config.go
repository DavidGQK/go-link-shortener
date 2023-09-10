package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerURL    string
	ShortURLBase string
}

var AppConfig Config

func GetConfig() *Config {
	flag.StringVar(&AppConfig.ServerURL, "a", "localhost:8080", "url where server runs on")
	flag.StringVar(&AppConfig.ShortURLBase, "b", "http://localhost:8080", "base url for shortened link")
	flag.Parse()

	if envServerURL := os.Getenv("SERVER_ADDRESS"); envServerURL != "" {
		AppConfig.ServerURL = envServerURL
	}

	if envServerURLBase := os.Getenv("BASE_URL"); envServerURLBase != "" {
		AppConfig.ShortURLBase = envServerURLBase
	}

	return &AppConfig
}
