package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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

func (s *SQLiteStorage) Store(ctx context.Context, entry *domain.ClipboardEntry) (*domain.ClipboardEntry, error) {
	result, err := s.db.ExecContext(ctx,
		"INSERT INTO clipboard_history (content, timestamp) VALUES (?, ?)",
		entry.Content, entry.Timestamp,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &domain.ClipboardEntry{
		Id:        strconv.FormatInt(id, 10),
		Content:   entry.Content,
		Timestamp: entry.Timestamp,
	}, nil
}

func (s *SQLiteStorage) GetRecent(ctx context.Context, n int) ([]*domain.ClipboardEntry, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, content, timestamp FROM clipboard_history ORDER BY timestamp DESC LIMIT ?",
		n,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*domain.ClipboardEntry

	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		var id int64
		var content string
		var timestamp time.Time
		if err := rows.Scan(&id, &content, &timestamp); err != nil {
			return nil, err
		}
		entries = append(entries, &domain.ClipboardEntry{
			Id:        strconv.FormatInt(id, 10),
			Content:   content,
			Timestamp: timestamp,
		})
	}

	return entries, nil
}

func (s *SQLiteStorage) GetById(ctx context.Context, id string) (*domain.ClipboardEntry, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	var content string
	var timestamp time.Time

	err = s.db.QueryRowContext(ctx,
		"SELECT content, timestamp FROM clipboard_history WHERE id = ?",
		idInt,
	).Scan(&content, &timestamp)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("entry not found")
	}
	if err != nil {
		return nil, err
	}

	return &domain.ClipboardEntry{
		Id:        id,
		Content:   content,
		Timestamp: timestamp,
	}, nil

}

func (s *SQLiteStorage) Delete(ctx context.Context, id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	result, err := s.db.ExecContext(ctx,
		"DELETE FROM clipboard_history WHERE id = ?",
		idInt,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("entry not found")
	}

	return nil
}

func (s *SQLiteStorage) Search(ctx context.Context, query string, limit int) ([]*domain.ClipboardEntry, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, content, timestamp FROM clipboard_history WHERE content LIKE ? ORDER BY timestamp DESC LIMIT ?",
		"%"+query+"%",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*domain.ClipboardEntry

	for rows.Next() {
		var id int64
		var content string
		var timestamp time.Time
		if err := rows.Scan(&id, &content, &timestamp); err != nil {
			return nil, err
		}
		entries = append(entries, &domain.ClipboardEntry{
			Id:        strconv.FormatInt(id, 10),
			Content:   content,
			Timestamp: timestamp,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (s *SQLiteStorage) Count(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM clipboard_history",
	).Scan(&count)

	return count, err
}

func (s *SQLiteStorage) Clear(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM clipboard_history")
	return err
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
