package service

import (
	"context"
	"fmt"

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

type Analyzer interface {
	Analyze(entry *domain.ClipboardEntry) *domain.Analysis
}

type ClipboardService struct {
	storage  Storage
	analyzer Analyzer
}

type Stats struct {
	TotalEntries int
}

func NewClipboardService(storage Storage, analyzer Analyzer) *ClipboardService {
	return &ClipboardService{
		storage:  storage,
		analyzer: analyzer,
	}
}

func (s *ClipboardService) ProcessNewEntry(ctx context.Context, entry *domain.ClipboardEntry) (*domain.ClipboardEntry, error) {
	if entry == nil {
		return nil, fmt.Errorf("entry cannot be nil")
	}

	if entry.Content == "" {
		return nil, ErrEmptyContent
	}

	analysis := s.analyzer.Analyze(entry)

	if analysis.IsSensitive {
		return nil, ErrSensitiveContent
	}

	stored, err := s.storage.Store(ctx, entry)
	if err != nil {
		return nil, fmt.Errorf("failed to store entry: %w", err)
	}

	return stored, nil
}

func (s *ClipboardService) GetHistory(ctx context.Context, limit int) ([]*domain.ClipboardEntry, error) {
	if limit <= 0 || limit > 100 {
		return nil, ErrInvalidLimit
	}

	entries, err := s.storage.GetRecent(ctx, limit)

	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}

	return entries, nil
}

func (s *ClipboardService) GetEntry(ctx context.Context, id string) (*domain.ClipboardEntry, error) {
	if id == "" {
		return nil, ErrInvalidId
	}

	entry, err := s.storage.GetById(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}

	return entry, nil
}

func (s *ClipboardService) DeleteEntry(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidId
	}

	err := s.storage.Delete(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	return nil
}

func (s *ClipboardService) Search(ctx context.Context, query string, limit int) ([]*domain.ClipboardEntry, error) {
	if query == "" {
		return nil, ErrEmptyQuery
	}

	if limit <= 0 || limit > 1000 {
		limit = 100 // default
	}

	entries, err := s.storage.Search(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return entries, nil
}

func (s *ClipboardService) ClearHistory(ctx context.Context) error {
	return s.storage.Clear(ctx)
}

func (s *ClipboardService) GetStats(ctx context.Context) (*Stats, error) {
	count, err := s.storage.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return &Stats{
		TotalEntries: count,
	}, nil
}
