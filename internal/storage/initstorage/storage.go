package initstorage

import (
	"context"
	"github.com/DavidGQK/go-link-shortener/internal/models"
	"github.com/DavidGQK/go-link-shortener/internal/storage/cachestorage"
	db "github.com/DavidGQK/go-link-shortener/internal/storage/db"
	"github.com/DavidGQK/go-link-shortener/internal/storage/filestorage"
)

const (
	MemoryMode = iota
	FileMode
	DBMode
)

type Storage struct {
	storage models.StorageInterface
}

func NewStorage(filename string, dbConnData string) (*Storage, error) {
	mode := MemoryMode

	if dbConnData != "" {
		mode = DBMode
		datab, err := db.NewDB(dbConnData, mode)
		if err != nil {
			return nil, err
		}

		return &Storage{storage: datab}, nil

	} else if filename != "" {
		mode = FileMode

		fStore, err := filestorage.NewFStor(filename, mode)
		if err != nil {
			return nil, err
		}

		return &Storage{storage: fStore}, nil
	}

	cacheStore, err := cachestorage.NewCacheStor(mode)
	if err != nil {
		return nil, err
	}
	return &Storage{storage: cacheStore}, nil
}

func (s *Storage) Restore() error {
	return s.storage.Restore()
}

func (s *Storage) Add(key, value, cookie string) error {
	return s.storage.Add(key, value, cookie)
}

func (s *Storage) AddBatch(ctx context.Context, records []models.Record) error {
	return s.storage.AddBatch(ctx, records)
}

func (s *Storage) Get(key string) (string, bool) {
	return s.storage.Get(key)
}

func (s *Storage) GetMode() int {
	return s.storage.GetMode()
}

func (s *Storage) GetByOriginURL(originURL string) (string, error) {
	return s.storage.GetByOriginURL(originURL)
}

func (s *Storage) HealthCheck() error {
	return s.storage.HealthCheck()
}

func (s *Storage) GetUserRecords(ctx context.Context, cookie string) ([]models.Record, error) {
	return s.storage.GetUserRecords(ctx, cookie)
}

func (s *Storage) FindUserByID(ctx context.Context, userID int) (*models.User, error) {
	return s.storage.FindUserByID(ctx, userID)
}

func (s *Storage) CreateUser(ctx context.Context) (*models.User, error) {
	return s.storage.CreateUser(ctx)
}

func (s *Storage) UpdateUser(ctx context.Context, id int, cookie string) error {
	return s.storage.UpdateUser(ctx, id, cookie)
}
