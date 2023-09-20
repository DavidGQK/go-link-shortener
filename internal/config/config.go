package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerURL    string
	ShortURLBase string
	LoggingLevel string
}

func loadFlagConfig(AppConfig *Config) {
	flag.StringVar(&AppConfig.ServerURL, "a", "localhost:8080", "url where server runs on")
	flag.StringVar(&AppConfig.ShortURLBase, "b", "http://localhost:8080", "base url for shortened link")
	flag.Parse()
}

func loadEnvConfig(AppConfig *Config) {
	if envServerURL := os.Getenv("SERVER_ADDRESS"); envServerURL != "" {
		AppConfig.ServerURL = envServerURL
	}

	if envServerURLBase := os.Getenv("BASE_URL"); envServerURLBase != "" {
		AppConfig.ShortURLBase = envServerURLBase
	}
}

func loadLogLevel(AppConfig *Config) {
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		AppConfig.LoggingLevel = envLogLevel
	}
}

func GetConfig() *Config {
	var AppConfig Config

	loadFlagConfig(&AppConfig)
	loadEnvConfig(&AppConfig)
	loadLogLevel(&AppConfig)

	return &AppConfig
}
