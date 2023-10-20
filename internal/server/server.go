package server

import (
	"context"
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/DavidGQK/go-link-shortener/internal/models"
)

type repository interface {
	Add(string, string, string) error
	Get(string) (string, bool)
	HealthCheck() error
	GetMode() int
	AddBatch(context.Context, []models.Record) error
	GetByOriginURL(string) (string, error)
	GetUserRecords(context.Context, string) ([]models.Record, error)
	FindUserByID(context.Context, int) (*models.User, error)
	CreateUser(context.Context) (*models.User, error)
	UpdateUser(context.Context, int, string) error
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
