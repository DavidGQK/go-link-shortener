package storage

import (
	"bufio"
	"encoding/json"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
)

type Record struct {
	UUID        uuid.UUID `json:"UUID"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

type dataWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func (p *dataWriter) WriteData(rec *Record) error {
	return p.encoder.Encode(rec)
}

func (p *dataWriter) Close() error {
	return p.file.Close()
}

func NewDataWriter(filename string) (*dataWriter, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Log.Error("open file error", zap.Error(err))
		return nil, err
	}

	return &dataWriter{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

type Storage struct {
	dataWriter *dataWriter
	filename   string
	links      map[string]string
}

func (s *Storage) Restore() error {
	file, err := os.OpenFile(s.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		var rec Record
		line := fileScanner.Text()
		err = json.Unmarshal([]byte(line), &rec)
		if err != nil {
			logger.Log.Error("data decoding error", zap.Error(err))
			continue
		}

		s.links[rec.ShortURL] = rec.OriginalURL
	}

	return nil
}

func (s *Storage) Add(key, value string) {
	id := uuid.New()
	rec := Record{
		UUID:        id,
		ShortURL:    key,
		OriginalURL: value,
	}

	if s.filename != "" {
		dataWr, err := NewDataWriter(s.filename)
		if err != nil {
			logger.Log.Error("open file error", zap.Error(err))
		}

		s.dataWriter = dataWr
		defer s.dataWriter.Close()

		err = s.dataWriter.WriteData(&rec)
		if err != nil {
			logger.Log.Error("error while writing data", zap.Error(err))
		}
	}

	s.links[key] = value
}

func (s *Storage) Get(key string) (string, bool) {
	value, found := s.links[key]
	return value, found
}

func New(filename string) (*Storage, error) {
	var dataWr *dataWriter

	store := Storage{
		dataWriter: dataWr,
		filename:   filename,
		links:      make(map[string]string),
	}

	return &store, nil
}
