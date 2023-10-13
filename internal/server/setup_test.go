package server

import (
	"context"
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/DavidGQK/go-link-shortener/internal/models"
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

func (s *TestStorage) Restore() error {
	return nil
}

func (s *TestStorage) Add(key, value string) error {
	s.links[key] = value
	return nil
}

func (s *TestStorage) AddBatch(ctx context.Context, _ []models.Record) error {
	return nil
}

func (s *TestStorage) Get(key string) (string, bool) {
	value, found := s.links[key]
	return value, found
}

func (s *TestStorage) GetMode() int {
	return 0
}

func (s *TestStorage) GetByOriginURL(_ string) (string, error) {
	return "", nil
}

func (s *TestStorage) HealthCheck() error {
	return nil
}
