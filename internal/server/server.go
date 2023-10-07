package server

import "github.com/DavidGQK/go-link-shortener/internal/config"

type repository interface {
	Add(string, string)
	Get(string) (string, bool)
}

type database interface {
	HealthCheck() error
}

type Server struct {
	config  *config.Config
	storage repository
	db      database
}

func New(c *config.Config, s repository) Server {
	return Server{
		config:  c,
		storage: s,
	}
}

func (s *Server) SetDB(db database) {
	s.db = db
}
