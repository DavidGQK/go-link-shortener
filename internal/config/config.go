package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerURL    string
	ShortURLBase string
	LoggingLevel string
	Filename     string
}

func loadFlagConfig(AppConfig *Config) {
	flag.StringVar(&AppConfig.ServerURL, "a", "localhost:8080", "url where server runs on")
	flag.StringVar(&AppConfig.ShortURLBase, "b", "http://localhost:8080", "base url for shortened link")
	flag.StringVar(&AppConfig.LoggingLevel, "l", "info", "logging level")
	flag.StringVar(&AppConfig.Filename, "f", "/tmp/short-url-db.json", "storage")

	flag.Parse()
}

func loadEnvConfig(AppConfig *Config) {
	if envServerURL := os.Getenv("SERVER_ADDRESS"); envServerURL != "" {
		AppConfig.ServerURL = envServerURL
	}

	if envServerURLBase := os.Getenv("BASE_URL"); envServerURLBase != "" {
		AppConfig.ShortURLBase = envServerURLBase
	}

	if envLoggingLevel := os.Getenv("LOG_LEVEL"); envLoggingLevel != "" {
		AppConfig.LoggingLevel = envLoggingLevel
	}

	if envFilename := os.Getenv("FILE_STORAGE_PATH"); envFilename != "" {
		AppConfig.Filename = envFilename
	}
}

func GetConfig() *Config {
	var AppConfig Config

	loadFlagConfig(&AppConfig)
	loadEnvConfig(&AppConfig)

	return &AppConfig
}
