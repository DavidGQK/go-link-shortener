package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/DavidGQK/go-link-shortener/internal/models"
	"github.com/DavidGQK/go-link-shortener/internal/storage/db"
	"github.com/DavidGQK/go-link-shortener/internal/storage/initstorage"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

func (s *Server) PostShortenLink(w http.ResponseWriter, r *http.Request) {
	var resp []byte
	var respStatus int

	initialURL, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	longURL, err := url.ParseRequestURI(string(initialURL))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	longURLStr := strings.Replace(longURL.String(), "%20", "", -1)
	if utf8.RuneCountInString(longURLStr) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := makeRandStringBytes(shortenedURLLength)
	err = s.storage.Add(id, longURLStr)
	if err != nil {
		if err == db.ErrConflict {
			id, err = s.storage.GetByOriginURL(longURLStr)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			respStatus = http.StatusConflict
			shortURLStr := fmt.Sprintf("%s/%s", s.config.ShortURLBase, id)
			resp = []byte(shortURLStr)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		respStatus = http.StatusCreated
		shortURLStr := fmt.Sprintf("%s/%s", s.config.ShortURLBase, id)
		resp = []byte(shortURLStr)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(respStatus)
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (s *Server) GetContent(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Path

	id := strings.Split(relPath, "/")[1]
	if utf8.RuneCountInString(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if longURLStr, ok := s.storage.Get(id); ok {
		w.Header().Set("Location", longURLStr)
		w.WriteHeader(http.StatusTemporaryRedirect)
		_, err := w.Write([]byte(longURLStr))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (s *Server) PostAPIShortenLink(w http.ResponseWriter, r *http.Request) {
	var body models.RequestShortenLink
	var resp models.ResponseShortenLink
	var respStatus int

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	longURLStr := strings.Replace(body.URL, " ", "", -1)
	if utf8.RuneCountInString(longURLStr) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := makeRandStringBytes(shortenedURLLength)
	err := s.storage.Add(id, longURLStr)
	if err != nil {
		if err == db.ErrConflict {
			id, err = s.storage.GetByOriginURL(longURLStr)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			respStatus = http.StatusConflict
			shortURLStr := fmt.Sprintf("%s/%s", s.config.ShortURLBase, id)
			resp = models.ResponseShortenLink{
				Result: shortURLStr,
			}
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		respStatus = http.StatusCreated
		shortURLStr := fmt.Sprintf("%s/%s", s.config.ShortURLBase, id)
		resp = models.ResponseShortenLink{
			Result: shortURLStr,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(respStatus)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	if s.storage.GetMode() != initstorage.DBMode {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err := s.storage.HealthCheck()
	if err != nil {
		logger.Log.Error(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Connection to DB is successful"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (s *Server) PostAPIShortenBatch(w http.ResponseWriter, r *http.Request) {
	var body models.RequestBatchLinks
	var records []models.Record
	var response models.ResponseBatchLinks

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, el := range body {
		id := makeRandStringBytes(shortenedURLLength)
		shortURLStr := s.config.ShortURLBase + "/" + id

		rec := models.Record{
			UUID:        el.CorrelationID,
			OriginalURL: el.OriginalURL,
			ShortURL:    id,
		}

		records = append(records, rec)

		response = append(response, models.ResponseLinks{
			CorrelationID: rec.UUID,
			ShortURL:      shortURLStr,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := s.storage.AddBatch(ctx, records)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const shortenedURLLength = 10

func makeRandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
