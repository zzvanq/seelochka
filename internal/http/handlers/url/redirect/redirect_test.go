package redirect_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zzvanq/seelochka/internal/http/handlers/url/redirect"
	"github.com/zzvanq/seelochka/internal/storage"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		alias      string
		url        string
		statusCode int
		mockError  error
	}{
		{
			alias:      "success",
			url:        "https://www.success.com/",
			statusCode: http.StatusFound,
			mockError:  nil,
		},
		{
			alias:      "fail",
			url:        "",
			statusCode: http.StatusNotFound,
			mockError:  storage.ErrURLNotFound,
		},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	for _, tc := range cases {
		t.Run(tc.alias, func(t *testing.T) {
			urlGetterMock := redirect.NewMockURLGetter(t)
			urlGetterMock.On("GetURL", tc.alias).Return(tc.url, tc.mockError).Once()

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(logger, urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			client := &http.Client{
				CheckRedirect: func(r *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}
			resp, err := client.Get(ts.URL + "/" + tc.alias)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tc.statusCode, resp.StatusCode)
			assert.Equal(t, tc.url, resp.Header.Get("Location"))
		})
	}
}
