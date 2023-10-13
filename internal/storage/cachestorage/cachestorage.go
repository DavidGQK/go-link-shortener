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

func (s *CacheStor) Add(key, value string) error {
	s.links[key] = value
	return nil
}

func (s *CacheStor) AddBatch(ctx context.Context, records []models.Record) error {
	for _, rec := range records {
		s.links[rec.ShortURL] = rec.OriginalURL
	}
	return nil
}

func (s *CacheStor) Get(key string) (string, bool) {
	value, found := s.links[key]
	return value, found
}

func (s *CacheStor) GetMode() int {
	return s.mode
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
