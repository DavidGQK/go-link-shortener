package server

import "github.com/DavidGQK/go-link-shortener/internal/config"

type repository interface {
	Add(string, string)
	Get(string) (string, bool)
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
