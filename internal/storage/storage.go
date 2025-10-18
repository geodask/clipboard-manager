package storage

import (
	"context"

	"github.com/geodask/clipboard-manager/internal/domain"
)

type Storage interface {
	Store(ctx context.Context, entry *domain.ClipboardEntry) (*domain.ClipboardEntry, error)
	GetRecent(ctx context.Context, n int) ([]*domain.ClipboardEntry, error)

	GetById(ctx context.Context, id string) (*domain.ClipboardEntry, error)
	Delete(ctx context.Context, id string) error

	Search(ctx context.Context, query string, limit int) ([]*domain.ClipboardEntry, error)
	Count(ctx context.Context) (int, error)

	Clear(ctx context.Context) error
}
