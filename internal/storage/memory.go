package storage

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/geodask/clipboard-manager/internal/domain"
)

type MemoryStorage struct {
	entries []*domain.ClipboardEntry
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (ms *MemoryStorage) Store(ctx context.Context, entry *domain.ClipboardEntry) (*domain.ClipboardEntry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	id := strconv.Itoa(len(ms.entries) + 1)
	storedEntry := &domain.ClipboardEntry{
		Id:        id,
		Content:   entry.Content,
		Timestamp: entry.Timestamp,
	}
	ms.entries = append(ms.entries, storedEntry)
	return storedEntry, nil
}

func (ms *MemoryStorage) GetRecent(ctx context.Context, n int) ([]*domain.ClipboardEntry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if len(ms.entries) < n {
		return ms.entries, nil
	}
	return ms.entries[len(ms.entries)-n:], nil
}

func (ms *MemoryStorage) GetByID(ctx context.Context, id string) (*domain.ClipboardEntry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	for _, entry := range ms.entries {
		if entry.Id == id {
			return entry, nil
		}
	}
	return nil, fmt.Errorf("entry not found")
}

func (ms *MemoryStorage) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	for i, entry := range ms.entries {
		if entry.Id == id {
			ms.entries = append(ms.entries[:i], ms.entries[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("entry not found")
}

func (ms *MemoryStorage) Search(ctx context.Context, query string, limit int) ([]*domain.ClipboardEntry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var results []*domain.ClipboardEntry

	for i := len(ms.entries) - 1; i >= 0 && len(results) < limit; i-- {
		if contains(ms.entries[i].Content, query) {
			results = append(results, ms.entries[i])
		}
	}

	return results, nil
}

func (ms *MemoryStorage) Count(ctx context.Context) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	return len(ms.entries), nil
}

func (ms *MemoryStorage) Clear(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	ms.entries = nil
	return nil
}

func contains(content, query string) bool {
	return strings.Contains(strings.ToLower(content), strings.ToLower(query))
}
