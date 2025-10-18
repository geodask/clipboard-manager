package storage

import (
	"strconv"

	"github.com/geodask/clipboard-manager/internal/domain"
)

type Storage interface {
	Store(entry *domain.ClipboardEntry) (*domain.ClipboardEntry, error)
	GetRecent(n int) ([]*domain.ClipboardEntry, error)
}

type MemoryStorage struct {
	entries []*domain.ClipboardEntry
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (ms *MemoryStorage) Store(entry *domain.ClipboardEntry) (*domain.ClipboardEntry, error) {
	id := strconv.Itoa(len(ms.entries) + 1)
	storedEntry := &domain.ClipboardEntry{
		Id:        id,
		Content:   entry.Content,
		Timestamp: entry.Timestamp,
	}
	ms.entries = append(ms.entries, storedEntry)
	return storedEntry, nil
}

func (ms *MemoryStorage) GetRecent(n int) ([]*domain.ClipboardEntry, error) {
	if len(ms.entries) < n {
		return ms.entries, nil
	}
	return ms.entries[len(ms.entries)-n:], nil
}
