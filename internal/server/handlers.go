package server

import (
	"encoding/json"
	"github.com/DavidGQK/go-link-shortener/internal/logger"
	"github.com/DavidGQK/go-link-shortener/internal/models"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

func (s *Server) PostShortenLink(w http.ResponseWriter, r *http.Request) {
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
	shortURLStr := s.config.ShortURLBase + "/" + id
	s.storage.Add(id, longURLStr)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURLStr))
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
	shortURLStr := s.config.ShortURLBase + "/" + id
	s.storage.Add(id, longURLStr)

	resp := models.ResponseShortenLink{
		Result: shortURLStr,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err := s.db.HealthCheck()
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const shortenedURLLength = 10

func makeRandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
