package save_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zzvanq/seelochka/internal/http/handlers/url/save"
	"github.com/zzvanq/seelochka/internal/storage"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name       string
		payload    map[string]string
		respError  string
		statusCode int
		mockError  error
	}{
		{
			name: "No url",
			payload: map[string]string{
				"url":   "",
				"alias": "test",
			},
			respError:  "field 'URL' is 'required'",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Not url",
			payload: map[string]string{
				"url":   "noturl@gmail.com",
				"alias": "test",
			},
			respError:  "field 'URL' is 'url'",
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Alias Used",
			payload: map[string]string{
				"url":   "https://used.com",
				"alias": "used",
			},
			respError:  "alias is used",
			mockError:  storage.ErrAliasUsed,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Random alias",
			payload: map[string]string{
				"url": "https://url.com",
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "Given alias",
			payload: map[string]string{
				"url":   "https://url.com",
				"alias": "alias",
			},
			statusCode: http.StatusCreated,
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := save.NewMockURLSaver(t)

			if tc.mockError != nil || tc.respError == "" {
				urlSaverMock.On("SaveURL",
					tc.payload["url"], mock.AnythingOfType("string")).Return(tc.mockError).Once()
			}
			handler := save.New(logger, urlSaverMock)

			rr := httptest.NewRecorder()
			payload, err := json.Marshal(tc.payload)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(payload)))
			require.NoError(t, err)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, tc.statusCode)

			body := rr.Body.String()
			var resp save.AliasResponse
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Error)

			if _, exists := tc.payload["alias"]; exists && tc.statusCode == http.StatusCreated {
				require.Equal(t, tc.payload["alias"], resp.Alias)
			}
		})
	}
}
