package storage

import (
	"database/sql"
	"time"

	"github.com/geodask/clipboard-manager/internal/domain"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS clipboard_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			timestamp DATETIME NOT NULL
		)
	`)

	if err != nil {
		db.Close()
		return nil, err
	}

	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) Store(entry *domain.ClipboardEntry) error {
	_, err := s.db.Exec(
		"INSERT INTO clipboard_history (content, timestamp) VALUES (?, ?)",
		entry.Content, entry.Timestamp,
	)
	return err
}

func (s *SQLiteStorage) GetRecent(n int) ([]*domain.ClipboardEntry, error) {
	rows, err := s.db.Query(
		"SELECT content, timestamp FROM clipboard_history ORDER BY timestamp DESC LIMIT ?",
		n,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*domain.ClipboardEntry

	for rows.Next() {
		var content string
		var timestamp time.Time
		if err := rows.Scan(&content, &timestamp); err != nil {
			return nil, err
		}
		entries = append(entries, &domain.ClipboardEntry{
			Content:   content,
			Timestamp: timestamp,
		})
	}

	return entries, nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
