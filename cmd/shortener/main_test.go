package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_processPOST(t *testing.T) {
	type args struct {
		body         string
		responseCode int
		contentType  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Response 201 - StatusCreated",
			args: args{
				body:         "https://practicum.yandex.ru/",
				responseCode: http.StatusCreated,
			},
		},

		{
			name: "Response 400 - StatusBadRequest",
			args: args{
				body:         " ",
				responseCode: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := strings.NewReader(tt.args.body)
			req := httptest.NewRequest(http.MethodPost, "/", reqBody)
			req.Header.Set("Content-Type", "text/plain")
			w := httptest.NewRecorder()

			processPOST(w, req)
			result := w.Result()
			defer result.Body.Close()

			if tt.args.responseCode == http.StatusCreated {
				assert.Equal(t, tt.args.responseCode, result.StatusCode)

				resultBody, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				assert.NotEmpty(t, string(resultBody))
			} else {
				assert.Equal(t, tt.args.responseCode, result.StatusCode)
			}
		})
	}
}

func Test_processGET(t *testing.T) {
	tests := []struct {
		name         string
		responseCode int
	}{
		{
			name:         "Response 400 - StatusBadRequest",
			responseCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/abcdf12345", nil)
			w := httptest.NewRecorder()

			processGET(w, req)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.responseCode, result.StatusCode)
		})
	}
}
