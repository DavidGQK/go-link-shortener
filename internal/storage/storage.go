package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
)

type record struct {
	uuid        uuid.UUID `json:"uuid"`
	shortURL    string    `json:"short_url"`
	originalURL string    `json:"original_url"`
}

type dataWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func (p *dataWriter) WriteData(rec *record) error {
	return p.encoder.Encode(rec)
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
	links      map[string]string
}

func (s *Storage) Restore() error {
	file, err := os.OpenFile(s.filename, os.O_RDONLY|os.O_CREATE, 0777)
	defer file.Close()
	if err != nil {
		return err
	}

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		var rec record
		line := fileScanner.Text()
		err = json.Unmarshal([]byte(line), &rec)
		if err != nil {
			logger.Log.Error("data decoding error", zap.Error(err))
			continue
		}

		s.links[rec.shortURL] = rec.originalURL
	}

	return nil
}

func (s *Storage) Add(key, value string) {
	id := uuid.New()
	rec := record{
		uuid:        id,
		shortURL:    key,
		originalURL: value,
	}

	if s.filename != "" {
		dataWr, err := NewDataWriter(s.filename)
		if err != nil {
			logger.Log.Error("open file error", zap.Error(err))
		}

		s.dataWriter = dataWr
		defer s.dataWriter.Close()

		s.dataWriter.WriteData(&rec)
		fmt.Println("added record to the file")
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
