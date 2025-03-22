package api

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrInvalidStatusCode = errors.New("invalid status code")

func GetRedirect(aliasUrl string) (string, error) {
	resp, err := http.Get(aliasUrl)
	if err != nil {
		return "", err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%w, %d", ErrInvalidStatusCode, resp.StatusCode)
	}

	return resp.Header.Get("Location"), nil
}
