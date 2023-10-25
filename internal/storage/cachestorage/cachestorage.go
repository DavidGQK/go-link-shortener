package cachestorage

import (
	"context"
	"errors"
	"github.com/DavidGQK/go-link-shortener/internal/models"
)

type CacheStor struct {
	links map[string]string
	mode  int
}

func NewCacheStor(mode int) (*CacheStor, error) {
	newCacheStor := &CacheStor{
		mode:  mode,
		links: make(map[string]string),
	}

	return newCacheStor, nil
}

func (s *CacheStor) Add(key, value, _ string) error {
	s.links[key] = value
	return nil
}

func (s *CacheStor) AddBatch(_ context.Context, records []models.Record) error {
	for _, rec := range records {
		s.links[rec.ShortURL] = rec.OriginalURL
	}
	return nil
}

func (s *CacheStor) Get(key string) (string, error) {
	value, found := s.links[key]
	if !found {
		return "", errors.New("key not found")
	}
	return value, nil
}

func (s *CacheStor) GetMode() int {
	return s.mode
}

func (s *CacheStor) GetByOriginURL(_ string) (string, error) {
	return "", errors.New("not database mode")
}

func (s *CacheStor) HealthCheck() error {
	return errors.New("not database mode")
}

func (s *CacheStor) Restore() error {
	return errors.New("not filestorage or database mode")
}

func (s *CacheStor) CloseStorage() error {
	return nil
}

func (s *CacheStor) GetUserRecords(_ context.Context, _ string) ([]models.Record, error) {
	return nil, errors.New("not database mode")
}

func (s *CacheStor) FindUserByID(_ context.Context, _ int) (*models.User, error) {
	return nil, errors.New("not database mode")
}

func (s *CacheStor) CreateUser(_ context.Context) (*models.User, error) {
	return nil, errors.New("not database mode")
}

func (s *CacheStor) UpdateUser(_ context.Context, _ int, _ string) error {
	return errors.New("not database mode")
}

func (s *CacheStor) DeleteUserURLs(_ context.Context, _ models.DeletedURLMessage) error {
	return errors.New("not database mode")
}
