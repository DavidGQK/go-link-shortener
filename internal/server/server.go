package server

import (
	"context"
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/DavidGQK/go-link-shortener/internal/models"
	"time"
)

type repository interface {
	Add(string, string, string) error
	Get(string) (string, error)
	HealthCheck() error
	GetMode() int
	AddBatch(context.Context, []models.Record) error
	GetByOriginURL(string) (string, error)
	GetUserRecords(context.Context, string) ([]models.Record, error)
	FindUserByID(context.Context, int) (*models.User, error)
	CreateUser(context.Context) (*models.User, error)
	UpdateUser(context.Context, int, string) error
	DeleteUserURLs(context.Context, models.DeletedURLMessage) error
}

type Server struct {
	config          *config.Config
	storage         repository
	DeletedURLsChan chan models.DeletedURLMessage
}

func New(c *config.Config, s repository) Server {
	server := Server{
		config:          c,
		storage:         s,
		DeletedURLsChan: make(chan models.DeletedURLMessage, 10),
	}

	go server.deleteMessageBatch()

	return server
}

func (s *Server) deleteMessageBatch() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for {
		select {
		case msg := <-s.DeletedURLsChan:
			err := s.storage.DeleteUserURLs(ctx, msg)
			if err != nil {
				logger.Log.Error(err)
				continue
			}
		default:
			continue
		}
	}
}
