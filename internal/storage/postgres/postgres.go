package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/qeery8/url-shortener.git/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"

	var id int64

	if err := s.db.QueryRow(
		"INSERT INTO url (alias, url) VALUES ($1, $2) RETURNING id",
		alias,
		urlToSave,
	).Scan(&id); err != nil {
		if errors.Is(err, storage.ErrURLExists) {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	var url string

	if err := s.db.QueryRow(
		"SELECT url FROM url WHERE alias = $1",
		alias,
	).Scan(&url); err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}
