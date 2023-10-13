package server

import (
	"context"
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/DavidGQK/go-link-shortener/internal/models"
)

type repository interface {
	Add(string, string) error
	Get(string) (string, bool)
	HealthCheck() error
	GetMode() int
	AddBatch(context.Context, []models.Record) error
	GetByOriginURL(string) (string, error)
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
