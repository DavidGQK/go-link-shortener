package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"os"
	"time"
)

const (
	MemoryMode = iota
	FileMode
	DBMode
)

var ErrConflict = errors.New(`already exists`)

type DBInterface interface {
	FindRecord(ctx context.Context, value string) (Record, error)
	FindRecordByOriginURL(ctx context.Context, value string) (Record, error)
	HealthCheck() error
	Close() error
	CreateDBScheme() error
	SaveRecord(context.Context, *Record) error
	SaveRecordsBatch(context.Context, []Record) error
}

type Record struct {
	UUID        string `json:"UUID"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type DataWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func (p *DataWriter) WriteData(rec *Record) error {
	return p.encoder.Encode(rec)
}

func (p *DataWriter) Close() error {
	return p.file.Close()
}

func NewDataWriter(file *os.File) (*DataWriter, error) {
	return &DataWriter{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

type Storage struct {
	dataWriter *DataWriter
	filename   string
	links      map[string]string
	db         DBInterface
	mode       int
}

func (s *Storage) Restore() error {

	switch s.mode {
	case FileMode:
		fileScanner := bufio.NewScanner(s.dataWriter.file)
		for fileScanner.Scan() {
			var rec Record
			line := fileScanner.Text()
			err := json.Unmarshal([]byte(line), &rec)
			if err != nil {
				logger.Log.Error("data decoding error", zap.Error(err))
				continue
			}

			s.links[rec.ShortURL] = rec.OriginalURL
		}
	case DBMode:
		err := s.db.CreateDBScheme()
		if err != nil {
			logger.Log.Error("db restoring error", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *Storage) Add(key, value string) error {
	id := uuid.NewString()
	rec := Record{
		UUID:        id,
		ShortURL:    key,
		OriginalURL: value,
	}

	if s.mode == DBMode {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := s.db.SaveRecord(ctx, &rec)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return ErrConflict
			}

			logger.Log.Error("error while writing data to db", zap.Error(err))
			return err
		}
	} else if s.mode == FileMode {
		err := s.dataWriter.WriteData(&rec)
		if err != nil {
			logger.Log.Error("error while writing data", zap.Error(err))
			return err
		}
	}

	s.links[key] = value
	return nil
}

func (s *Storage) AddBatch(ctx context.Context, records []Record) error {
	if s.mode == DBMode {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := s.db.SaveRecordsBatch(ctx, records)
		if err != nil {
			logger.Log.Error("error while writing data batch to db", zap.Error(err))
		}
	}

	for _, rec := range records {
		if s.mode == FileMode {
			err := s.dataWriter.WriteData(&rec)
			if err != nil {
				logger.Log.Error("error while writing data in batch", zap.Error(err))
			}
		}

		s.links[rec.ShortURL] = rec.OriginalURL
	}

	return nil
}

func (s *Storage) Get(key string) (string, bool) {
	if s.mode == DBMode {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		rec, err := s.db.FindRecord(ctx, key)
		if err != nil {
			logger.Log.Error("db search query error", zap.Error(err))
			return "", false
		}

		return rec.OriginalURL, true
	}

	value, found := s.links[key]
	return value, found
}

func New(filename string, dbConnData string, dataWr *DataWriter) (*Storage, error) {
	var db *Database
	var err error

	mode := MemoryMode

	if dbConnData != "" {
		mode = DBMode
		db, err = NewDB(dbConnData)
		if err != nil {
			return nil, err
		}
	} else if filename != "" {
		mode = FileMode
	}

	store := Storage{
		dataWriter: dataWr,
		filename:   filename,
		links:      make(map[string]string),
		db:         db,
		mode:       mode,
	}

	return &store, nil
}

func (s *Storage) HealthCheck() error {
	return s.db.HealthCheck()
}

func (s *Storage) CloseConnection() error {
	return s.db.Close()
}

func (s *Storage) GetMode() int {
	return s.mode
}

func (s *Storage) GetByOriginURL(originURL string) (string, error) {
	if s.mode == DBMode {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		fmt.Println("ORIGIN URL = ", originURL)
		rec, err := s.db.FindRecordByOriginURL(ctx, originURL)
		if err != nil {
			return "", err
		}

		return rec.ShortURL, nil
	}

	return "", errors.New("not databas mode")
}
