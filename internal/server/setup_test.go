package server

import (
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/DavidGQK/go-link-shortener/internal/storage"
)

var TestCfg = config.Config{
	ServerURL:    "localhost:8080",
	ShortURLBase: "http://localhost:8080/",
}

type TestStorage struct {
	links map[string]string
}

func NewTestStorage() *TestStorage {
	return &TestStorage{
		links: make(map[string]string),
	}
}

func (s *TestStorage) Add(key, value string) {
	s.links[key] = value
}

func (s *TestStorage) Get(key string) (string, bool) {
	value, found := s.links[key]
	return value, found
}

func (s *TestStorage) HealthCheck() error {
	return nil
}

func (s *TestStorage) GetMode() int {
	return 0
}

func (s *TestStorage) AddBatch(_ []storage.Record) error {
	return nil
}
