package main

import (
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Maps for processed links
var shortToLong = make(map[string]string)
var longToShort = make(map[string]string)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const shortenedURLLength = 10
const baseURL = "http://localhost:8080/"

func handlerForShortening(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		processPOST(w, r)
	case http.MethodGet:
		processGET(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func processPOST(w http.ResponseWriter, r *http.Request) {
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

	if shortURLStr, ok := longToShort[longURLStr]; ok {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURLStr))
	} else {
		id := RandStringBytes(shortenedURLLength)
		shortURLStr = baseURL + id

		longToShort[longURLStr] = shortURLStr
		shortToLong[id] = longURLStr

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURLStr))
	}
}

func processGET(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Path

	id := strings.Split(relPath, "/")[1]
	if utf8.RuneCountInString(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if longURLStr, ok := shortToLong[id]; ok {
		w.Header().Set("Location", longURLStr)
		w.WriteHeader(http.StatusTemporaryRedirect)
		//w.WriteHeader(http.StatusCreated)
		w.Write([]byte(longURLStr))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func main() {
	r := chi.NewRouter()
	r.Get("/{id}", processGET)
	r.Post("/", processPOST)

	log.Fatal(http.ListenAndServe(":8080", r))
}
