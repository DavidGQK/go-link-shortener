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
	DataWriter *DataWriter
	Filename   string
	Links      map[string]string
	Mode       int
}

type DataWriter struct {
	File    *os.File
	Encoder *json.Encoder
}

func (p *DataWriter) WriteData(rec *models.Record) error {
	return p.Encoder.Encode(rec)
}

func (p *DataWriter) Close() error {
	return p.File.Close()
}

func NewDataWriter(file *os.File) (*DataWriter, error) {
	return &DataWriter{
		File:    file,
		Encoder: json.NewEncoder(file),
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
	//defer dataWr.Close()

	newFStor := &FStor{
		DataWriter: dataWr,
		Filename:   filename,
		Mode:       mode,
		Links:      make(map[string]string),
	}

	return newFStor, nil
}

func (s *FStor) Restore() error {
	fileScanner := bufio.NewScanner(s.DataWriter.File)
	for fileScanner.Scan() {
		var rec models.Record
		line := fileScanner.Text()
		err := json.Unmarshal([]byte(line), &rec)
		if err != nil {
			logger.Log.Error("data decoding error", zap.Error(err))
			continue
		}

		s.Links[rec.ShortURL] = rec.OriginalURL
	}

	return nil
}

func (s *FStor) Add(key, value string) error {
	id := uuid.NewString()
	rec := models.Record{
		UUID:        id,
		ShortURL:    key,
		OriginalURL: value,
	}

	err := s.DataWriter.WriteData(&rec)
	if err != nil {
		logger.Log.Error("error while writing data", zap.Error(err))
		return err
	}

	s.Links[key] = value
	return nil
}

func (s *FStor) AddBatch(ctx context.Context, records []models.Record) error {
	for _, rec := range records {
		err := s.DataWriter.WriteData(&rec)
		if err != nil {
			logger.Log.Error("error while writing data in batch", zap.Error(err))
		}

		s.Links[rec.ShortURL] = rec.OriginalURL
	}

	return nil
}

func (s *FStor) Get(key string) (string, bool) {
	value, found := s.Links[key]
	return value, found
}

func (s *FStor) GetMode() int {
	return s.Mode
}

func (s *FStor) GetByOriginURL(originURL string) (string, error) {
	return "", errors.New("not database mode")
}

func (s *FStor) HealthCheck() error {
	return errors.New("not database mode")
}

func (s *FStor) CloseStorage() error {
	return s.DataWriter.Close()
}
