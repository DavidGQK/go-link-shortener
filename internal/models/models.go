package models

import (
	"context"
	"errors"
)

var ErrConflict = errors.New(`already exists`)
var ErrDeleted = errors.New(`was deleted`)

type RequestShortenLink struct {
	URL string `json:"url"`
}

type ResponseShortenLink struct {
	Result string `json:"result"`
}

type RequestLinks struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type RequestBatchLinks []RequestLinks

type ResponseLinks struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type ResponseBatchLinks []ResponseLinks

type Record struct {
	UUID        string `json:"UUID"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	DeletedFlag bool   `json:"is_deleted"`
}

type StorageInterface interface {
	Restore() error
	Add(string, string, string) error
	AddBatch(context.Context, []Record) error
	Get(string) (string, error)
	GetMode() int
	GetByOriginURL(string) (string, error)
	HealthCheck() error
	CloseStorage() error
	GetUserRecords(context.Context, string) ([]Record, error)
	FindUserByID(context.Context, int) (*User, error)
	CreateUser(context.Context) (*User, error)
	UpdateUser(context.Context, int, string) error
}

type User struct {
	UserID int    `json:"user_id"`
	Cookie string `json:"cookie"`
}

type ResponseUserURLs []ResponseUserURL

type ResponseUserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
