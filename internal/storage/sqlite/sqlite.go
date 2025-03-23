package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/zzvanq/seelochka/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s Storage) SaveURL(longURL, alias string) error {
	stmt, err := s.db.Prepare("INSERT INTO urls(url, alias) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("prepare SaveURL: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(longURL, alias)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return storage.ErrAliasUsed
		}

		return fmt.Errorf("execute SaveURL: %w", err)
	}

	return nil
}

func (s Storage) GetURL(alias string) (string, error) {
	stmt, err := s.db.Prepare("SELECT url FROM urls WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("prepare GetURL: %w", err)
	}

	defer stmt.Close()

	var longURL string

	err = stmt.QueryRow(alias).Scan(&longURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("execute GetURL: %w", err)
	}

	return longURL, nil
}
