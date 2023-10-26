package server

import (
	"context"
	"errors"
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

func (s *TestStorage) Add(key, value, _ string) error {
	s.links[key] = value
	return nil
}

func (s *TestStorage) AddBatch(_ context.Context, _ []models.Record) error {
	return nil
}

func (s *TestStorage) Get(key string) (string, error) {
	value, found := s.links[key]
	if !found {
		return "", errors.New("key not found")
	}
	return value, nil
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

func (s *TestStorage) GetUserRecords(_ context.Context, _ string) ([]models.Record, error) {
	return nil, nil
}

func (s *TestStorage) FindUserByID(_ context.Context, _ int) (*models.User, error) {
	return nil, nil
}

func (s *TestStorage) CreateUser(_ context.Context) (*models.User, error) {
	return nil, nil
}

func (s *TestStorage) UpdateUser(_ context.Context, _ int, _ string) error {
	return nil
}

func (s *TestStorage) DeleteUserURLs(_ context.Context, _ models.DeletedURLMessage) error {
	return nil
}
