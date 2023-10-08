package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
	"reflect"
	"time"
)

const (
	MemoryMode = iota
	FileMode
	DBMode
)

type DBInterface interface {
	FindRecord(ctx context.Context, value string) (Record, error)
	HealthCheck() error
	Close() error
	Restore([]Record) error
	SaveRecord(context.Context, *Record) error
}

type Record struct {
	UUID        uuid.UUID `json:"UUID"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
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

func NewDataWriter(file *os.File, filename string) (*DataWriter, error) {
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

//func (s *Storage) Restore() error {
//	fileScanner := bufio.NewScanner(s.dataWriter.file)
//	for fileScanner.Scan() {
//		var rec Record
//		line := fileScanner.Text()
//		err := json.Unmarshal([]byte(line), &rec)
//		if err != nil {
//			logger.Log.Error("data decoding error", zap.Error(err))
//			continue
//		}
//
//		s.links[rec.ShortURL] = rec.OriginalURL
//	}
//
//	return nil
//}

func (s *Storage) Restore() error {
	var recs []Record
	fileScanner := bufio.NewScanner(s.dataWriter.file)
	for fileScanner.Scan() {
		var rec Record
		line := fileScanner.Text()
		err := json.Unmarshal([]byte(line), &rec)
		if err != nil {
			logger.Log.Error("data decoding error", zap.Error(err))
			continue
		}

		recs = append(recs, rec)
		s.links[rec.ShortURL] = rec.OriginalURL
	}

	if err := s.HealthCheck(); err == nil {
		err := s.db.Restore(recs)
		if err != nil {
			logger.Log.Error("db restoring error", zap.Error(err))
		}
	}

	return nil
}

//func (s *Storage) Add(key, value string) {
//	id := uuid.New()
//	rec := Record{
//		UUID:        id,
//		ShortURL:    key,
//		OriginalURL: value,
//	}
//
//	if s.filename != "" {
//		err := s.dataWriter.WriteData(&rec)
//		if err != nil {
//			logger.Log.Error("error while writing data", zap.Error(err))
//		}
//	}
//
//	s.links[key] = value
//}

func (s *Storage) Add(key, value string) {
	id := uuid.New()
	rec := Record{
		UUID:        id,
		ShortURL:    key,
		OriginalURL: value,
	}

	if err := s.HealthCheck(); err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := s.db.SaveRecord(ctx, &rec)
		if err != nil {
			logger.Log.Error("error while writing data to db", zap.Error(err))
		}
	} else if s.filename != "" {
		err := s.dataWriter.WriteData(&rec)
		if err != nil {
			logger.Log.Error("error while writing data", zap.Error(err))
		}
	}

	s.links[key] = value
}

//func (s *Storage) Get(key string) (string, bool) {
//	value, found := s.links[key]
//	return value, found
//}

func (s *Storage) Get(key string) (string, bool) {
	if err := s.HealthCheck(); err != nil {
		logger.Log.Error("db connection error", zap.Error(err))
		value, found := s.links[key]
		return value, found
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		rec, err := s.db.FindRecord(ctx, key)
		if err != nil {
			logger.Log.Error("db search query error", zap.Error(err))
			return "", false
		}

		return rec.OriginalURL, true
	}
}

func New(filename string, dbConnData string, dataWr *DataWriter) (*Storage, error) {
	var db *Database
	var err error

	if dbConnData != "" {
		db, err = NewDB(dbConnData)
		if err != nil {
			return nil, err
		}
	}

	store := Storage{
		dataWriter: dataWr,
		filename:   filename,
		links:      make(map[string]string),
		db:         db,
	}

	return &store, nil
}

func (s *Storage) HealthCheck() error {
	if reflect.ValueOf(s.db).IsNil() {
		return fmt.Errorf("nil db")
	}
	return s.db.HealthCheck()
}

func (s *Storage) CloseConnection() error {
	return s.db.Close()
}
