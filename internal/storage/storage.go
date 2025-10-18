package storage

import (
	"context"
	"strconv"

	"github.com/geodask/clipboard-manager/internal/domain"
)

type Storage interface {
	Store(ctx context.Context, entry *domain.ClipboardEntry) (*domain.ClipboardEntry, error)
	GetRecent(ctx context.Context, n int) ([]*domain.ClipboardEntry, error)
}

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
