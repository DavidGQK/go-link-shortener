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
}

func (s *Storage) Restore() error {
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
		err := s.dataWriter.WriteData(&rec)
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

func New(filename string, dataWr *DataWriter) (*Storage, error) {
	store := Storage{
		dataWriter: dataWr,
		filename:   filename,
		links:      make(map[string]string),
	}

	return &store, nil
}
