package models

import (
	"context"
)

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
}

type StorageInterface interface {
	Restore() error
	Add(string, string) error
	AddBatch(context.Context, []Record) error
	Get(string) (string, bool)
	GetMode() int
	GetByOriginURL(string) (string, error)
	HealthCheck() error
}
