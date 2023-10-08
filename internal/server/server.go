package server

import (
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/DavidGQK/go-link-shortener/internal/storage"
)

type repository interface {
	Add(string, string)
	Get(string) (string, bool)
	HealthCheck() error
	GetMode() int
	AddBatch(records []storage.Record) error
}

type Server struct {
	config  *config.Config
	storage repository
}

func New(c *config.Config, s repository) Server {
	return Server{
		config:  c,
		storage: s,
	}
}
