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
	UUID        uuid.UUID `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
}

type dataWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func (p *dataWriter) WriteData(record *Record) error {
	return p.encoder.Encode(record)
}

func (p *dataWriter) Close() error {
	return p.file.Close()
}

func NewDataWriter(filename string) (*dataWriter, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	defer file.Close()
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
	Links      map[string]string
}

func (s *Storage) Restore() error {
	file, err := os.OpenFile(s.filename, os.O_RDONLY|os.O_CREATE, 0777)
	defer file.Close()
	if err != nil {
		return err
	}

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		var record Record
		line := fileScanner.Text()
		err = json.Unmarshal([]byte(line), &record)
		if err != nil {
			logger.Log.Error("data decoding error", zap.Error(err))
			continue
		}

		s.Links[record.ShortURL] = record.OriginalURL
	}

	return nil
}

func (s *Storage) Add(key, value string) {
	id := uuid.New()
	record := Record{
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

		s.dataWriter.WriteData(&record)
	}

	s.Links[key] = value

}

func (s *Storage) Get(key string) (string, bool) {
	value, found := s.Links[key]
	return value, found
}

func New(filename string) (*Storage, error) {
	var dataWr *dataWriter

	store := Storage{
		dataWriter: dataWr,
		filename:   filename,
		Links:      make(map[string]string),
	}

	return &store, nil
}
