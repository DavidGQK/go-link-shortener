package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_PostShortLink(t *testing.T) {
	type fields struct {
		serverURL string
		storage   repository
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
				serverURL: TestCfg.ServerURL,
				storage:   NewTestStorage(),
			},
			want: want{
				expectedCode: http.StatusCreated,
			},
		},
		{
			name: "Response 400 - StatusBadRequest",
			body: " ",
			fields: fields{
				serverURL: TestCfg.ShortURLBase,
				storage:   NewTestStorage(),
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
				baseURL: tt.fields.serverURL,
				storage: tt.fields.storage,
			}

			s.PostShortLink(w, req)
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

func Test_GetContent(t *testing.T) {
	type fields struct {
		serverURL string
		baseURL   string
		storage   repository
		id        string
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
				serverURL: TestCfg.ServerURL,
				baseURL:   TestCfg.ShortURLBase,
				storage:   NewTestStorage(),
				id:        "abcdf12345",
			},
			want: want{
				expectedCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.fields.baseURL+tt.fields.id, nil)
			w := httptest.NewRecorder()

			s := Server{
				baseURL: tt.fields.serverURL,
				storage: tt.fields.storage,
			}

			s.GetContent(w, req)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.expectedCode, result.StatusCode)
		})
	}
}
