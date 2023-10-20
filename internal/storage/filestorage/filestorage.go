package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/DavidGQK/go-link-shortener/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
)

type FStor struct {
	dataWriter *DataWriter
	filename   string
	links      map[string]string
	mode       int
}

type DataWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func (p *DataWriter) WriteData(rec *models.Record) error {
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

func NewFStor(filename string, mode int) (*FStor, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Log.Error("open filestorage error", zap.Error(err))
		return nil, err
	}

	dataWr, err := NewDataWriter(file)
	if err != nil {
		logger.Log.Error("creating a new data writer error", zap.Error(err))
		return nil, err
	}

	newFStor := &FStor{
		dataWriter: dataWr,
		filename:   filename,
		mode:       mode,
		links:      make(map[string]string),
	}

	return newFStor, nil
}

func (s *FStor) Restore() error {
	fileScanner := bufio.NewScanner(s.dataWriter.file)
	for fileScanner.Scan() {
		var rec models.Record
		line := fileScanner.Text()
		err := json.Unmarshal([]byte(line), &rec)
		if err != nil {
			logger.Log.Error("data decoding error", zap.Error(err))
			continue
		}

		s.links[rec.ShortURL] = rec.OriginalURL
	}

	return nil
}

func (s *FStor) Add(key, value, _ string) error {
	id := uuid.NewString()
	rec := models.Record{
		UUID:        id,
		ShortURL:    key,
		OriginalURL: value,
	}

	err := s.dataWriter.WriteData(&rec)
	if err != nil {
		logger.Log.Error("error while writing data", zap.Error(err))
		return err
	}

	s.links[key] = value
	return nil
}

func (s *FStor) AddBatch(_ context.Context, records []models.Record) error {
	for _, rec := range records {
		err := s.dataWriter.WriteData(&rec)
		if err != nil {
			logger.Log.Error("error while writing data in batch", zap.Error(err))
		}

		s.links[rec.ShortURL] = rec.OriginalURL
	}

	return nil
}

func (s *FStor) Get(key string) (string, bool) {
	value, found := s.links[key]
	return value, found
}

func (s *FStor) GetMode() int {
	return s.mode
}

func (s *FStor) GetByOriginURL(_ string) (string, error) {
	return "", errors.New("not database mode")
}

func (s *FStor) HealthCheck() error {
	return errors.New("not database mode")
}

func (s *FStor) CloseStorage() error {
	return s.dataWriter.Close()
}

func (s *FStor) GetUserRecords(ctx context.Context, cookie string) ([]models.Record, error) {
	return nil, errors.New("not database mode")
}

func (s *FStor) FindUserByID(ctx context.Context, userID int) (*models.User, error) {
	return nil, errors.New("not database mode")
}

func (s *FStor) CreateUser(ctx context.Context) (*models.User, error) {
	return nil, errors.New("not database mode")
}

func (s *FStor) UpdateUser(ctx context.Context, id int, cookie string) error {
	return errors.New("not database mode")
}
