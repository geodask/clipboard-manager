package storage

import (
	"github.com/geodask/clipboard-manager/internal/domain"
)

type Storage interface {
	Store(entry *domain.ClipboardEntry) error
	GetRecent(n int) ([]*domain.ClipboardEntry, error)
}

type MemoryStorage struct {
	entries []*domain.ClipboardEntry
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (ms *MemoryStorage) Store(entry *domain.ClipboardEntry) error {
	ms.entries = append(ms.entries, entry)
	return nil
}

func (ms *MemoryStorage) GetRecent(n int) ([]*domain.ClipboardEntry, error) {
	if len(ms.entries) < n {
		return ms.entries, nil
	}
	return ms.entries[len(ms.entries)-n:], nil
}
