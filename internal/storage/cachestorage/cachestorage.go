package cachestorage

import (
	"context"
	"errors"
	"github.com/DavidGQK/go-link-shortener/internal/models"
)

type CacheStor struct {
	Links map[string]string
	Mode  int
}

func NewCacheStor(mode int) (*CacheStor, error) {
	newCacheStor := &CacheStor{
		Mode:  mode,
		Links: make(map[string]string),
	}

	return newCacheStor, nil
}

func (s *CacheStor) Add(key, value string) error {
	s.Links[key] = value
	return nil
}

func (s *CacheStor) AddBatch(ctx context.Context, records []models.Record) error {
	for _, rec := range records {
		s.Links[rec.ShortURL] = rec.OriginalURL
	}
	return nil
}

func (s *CacheStor) Get(key string) (string, bool) {
	value, found := s.Links[key]
	return value, found
}

func (s *CacheStor) GetMode() int {
	return s.Mode
}

func (s *CacheStor) GetByOriginURL(originURL string) (string, error) {
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
