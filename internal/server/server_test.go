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

func Test_PostShortenLink(t *testing.T) {
	type fields struct {
		config  *config.Config
		storage repository
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
				config:  &TestCfg,
				storage: NewTestStorage(),
			},
			want: want{
				expectedCode: http.StatusCreated,
			},
		},
		{
			name: "Response 400 - StatusBadRequest",
			body: " ",
			fields: fields{
				config:  &TestCfg,
				storage: NewTestStorage(),
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
				config:  tt.fields.config,
				storage: tt.fields.storage,
			}

			req.AddCookie(&http.Cookie{
				Name:  "shortener_session",
				Value: "test",
			})

			s.PostShortenLink(w, req)
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
		config  *config.Config
		storage repository
		id      string
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
				config:  &TestCfg,
				storage: NewTestStorage(),
				id:      "abcdf12345",
			},
			want: want{
				expectedCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.fields.config.ShortURLBase+tt.fields.id, nil)
			w := httptest.NewRecorder()

			s := Server{
				config:  tt.fields.config,
				storage: tt.fields.storage,
			}

			s.GetContent(w, req)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.expectedCode, result.StatusCode)
		})
	}
}

func Test_PostAPIShortenLink(t *testing.T) {
	type fields struct {
		config  *config.Config
		storage repository
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
			name: "POST API Response 201 - StatusCreated",
			body: `{ "url": "https://practicum.yandex.ru/" }`,
			fields: fields{
				config:  &TestCfg,
				storage: NewTestStorage(),
			},
			want: want{
				expectedCode: http.StatusCreated,
			},
		},
		{
			name: "Response 400 - StatusBadRequest",
			body: `{ "url": " " }`,
			fields: fields{
				config:  &TestCfg,
				storage: NewTestStorage(),
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
			req.Header.Set("Content-Type", "application/json")

			s := Server{
				config:  tt.fields.config,
				storage: tt.fields.storage,
			}

			req.AddCookie(&http.Cookie{
				Name:  "shortener_session",
				Value: "test",
			})

			s.PostAPIShortenLink(w, req)
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

func Test_PostAPIShortenBatch(t *testing.T) {
	type fields struct {
		config  *config.Config
		storage repository
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
			name: "POST API Response 201 - StatusCreated",
			body: `[
						{
							"correlation_id": "12345",
							"original_url": "https://practicum.yandex.ru/1"
						},
						{
							"correlation_id": "54321",
							"original_url": "https://practicum.yandex.ru/2"
						}
					]`,
			fields: fields{
				config:  &TestCfg,
				storage: NewTestStorage(),
			},
			want: want{
				expectedCode: http.StatusCreated,
			},
		},
		{
			name: "Response 400 - StatusBadRequest",
			body: ``,
			fields: fields{
				config:  &TestCfg,
				storage: NewTestStorage(),
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
			req.Header.Set("Content-Type", "application/json")

			s := Server{
				config:  tt.fields.config,
				storage: tt.fields.storage,
			}

			s.PostAPIShortenBatch(w, req)
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
