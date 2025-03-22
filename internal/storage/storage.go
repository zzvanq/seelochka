package storage

import (
	"errors"
)

var (
	ErrURLNotFound error = errors.New("url not found")
	ErrAliasUsed   error = errors.New("alias is used")
)
