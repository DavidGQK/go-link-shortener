package server

import (
	"encoding/json"
	"fmt"
	"github.com/DavidGQK/go-link-shortener/internal/models"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const shortenedURLLength = 10

type repository interface {
	AddToShort(string, string)
	GetFromShort(string) (string, bool)
	AddToLong(string, string)
	GetFromLong(string) (string, bool)
}

type Server struct {
	baseURL string
	storage repository
}

func NewServer(u string, s repository) Server {
	return Server{
		baseURL: u,
		storage: s,
	}
}

func (s Server) PostShortenLink(w http.ResponseWriter, r *http.Request) {
	initialURL, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	longURL, err := url.ParseRequestURI(string(initialURL))
	fmt.Println("longURL", longURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	longURLStr := strings.Replace(longURL.String(), "%20", "", -1)
	if utf8.RuneCountInString(longURLStr) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	if shortURLStr, ok := s.storage.GetFromLong(longURLStr); ok {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURLStr))
	} else {
		id := makeRandStringBytes(shortenedURLLength)
		shortURLStr = s.baseURL + "/" + id

		s.storage.AddToLong(longURLStr, shortURLStr)
		s.storage.AddToShort(id, longURLStr)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURLStr))
	}
}

func (s Server) GetContent(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Path

	id := strings.Split(relPath, "/")[1]
	if utf8.RuneCountInString(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if longURLStr, ok := s.storage.GetFromShort(id); ok {
		w.Header().Set("Location", longURLStr)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(longURLStr))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (s Server) PostAPIShortenLink(w http.ResponseWriter, r *http.Request) {
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

	shortURLStr, ok := s.storage.GetFromLong(longURLStr)
	if !ok {
		id := makeRandStringBytes(shortenedURLLength)
		shortURLStr = s.baseURL + "/" + id

		s.storage.AddToLong(longURLStr, shortURLStr)
		s.storage.AddToShort(id, longURLStr)
	}
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

func makeRandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
