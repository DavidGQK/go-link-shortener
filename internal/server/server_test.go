package server

import (
	"github.com/DavidGQK/go-link-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_ProcessPOST(t *testing.T) {
	type fields struct {
		Config  *config.Config
		Storage Repository
	}

	type want struct {
		expectedCode int
	}

	tests := []struct {
		name   string
		body   string
		fields fields
		want   want
	}{
		{
			name: "Response 201 - StatusCreated",
			body: "https://practicum.yandex.ru/",
			fields: fields{
				Config:  TestConfig,
				Storage: NewTestStorage(),
			},
			want: want{
				expectedCode: http.StatusCreated,
			},
		},
		{
			name: "Response 400 - StatusBadRequest",
			body: " ",
			fields: fields{
				Config:  TestConfig,
				Storage: NewTestStorage(),
			},
			want: want{
				expectedCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			reqBody := strings.NewReader(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/", reqBody)
			req.Header.Set("Content-Type", "text/plain")

			s := Server{
				Config:  tt.fields.Config,
				Storage: tt.fields.Storage,
			}

			s.ProcessPOST(w, req)
			result := w.Result()
			defer result.Body.Close()

			if tt.want.expectedCode == http.StatusCreated {
				assert.Equal(t, tt.want.expectedCode, result.StatusCode)

				resultBody, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				assert.NotEmpty(t, string(resultBody))
			} else {
				assert.Equal(t, tt.want.expectedCode, result.StatusCode)
			}
		})
	}
}

func Test_ProcessGET(t *testing.T) {
	type fields struct {
		Config  *config.Config
		Storage Repository
	}

	type want struct {
		expectedCode int
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "Response 400 - StatusBadRequest",
			fields: fields{
				Config:  TestConfig,
				Storage: NewTestStorage(),
			},
			want: want{
				expectedCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/abcdf12345", nil)
			w := httptest.NewRecorder()

			s := Server{
				Config:  tt.fields.Config,
				Storage: tt.fields.Storage,
			}

			s.ProcessGET(w, req)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.expectedCode, result.StatusCode)
		})
	}
}
