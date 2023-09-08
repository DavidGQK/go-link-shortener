package handlers

import (
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/DavidGQK/go-link-shortener/internal/storage"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const shortenedURLLength = 10

func handlerForShortening(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ProcessPOST(w, r)
	case http.MethodGet:
		ProcessGET(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func ProcessPOST(w http.ResponseWriter, r *http.Request) {
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

	if shortURLStr, ok := storage.LongToShort[longURLStr]; ok {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURLStr))
	} else {
		id := makeRandStringBytes(shortenedURLLength)
		shortURLStr = config.ShortURLBase + "/" + id

		storage.LongToShort[longURLStr] = shortURLStr
		storage.ShortToLong[id] = longURLStr

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURLStr))
	}
}

func ProcessGET(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Path

	id := strings.Split(relPath, "/")[1]
	if utf8.RuneCountInString(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if longURLStr, ok := storage.ShortToLong[id]; ok {
		w.Header().Set("Location", longURLStr)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(longURLStr))
	} else {
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
