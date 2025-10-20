package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/geodask/clipboard-manager/internal/domain"
)

type MockStorage struct {
	StoreResult           *domain.ClipboardEntry
	StoreError            error
	GetRecentResult       []*domain.ClipboardEntry
	GetRecentError        error
	GetByIdResult         *domain.ClipboardEntry
	GetByIdError          error
	DeleteError           error
	SearchResult          []*domain.ClipboardEntry
	SearchError           error
	CountResult           int
	CountError            error
	ClearError            error
	DeleteOlderThanResult int
	DeleteOlderThanError  error

	StoreCalled           bool
	StoreCalledWith       *domain.ClipboardEntry
	GetRecentCalled       bool
	GetRecentLimit        int
	GetByIdCalled         bool
	GetByIdId             string
	DeleteCalled          bool
	DeleteId              string
	SearchCalled          bool
	SearchQuery           string
	SearchLimit           int
	CountCalled           bool
	ClearCalled           bool
	DeleteOlderThanCalled bool
	DeleteOlderThanCutoff time.Time
}

func (m *MockStorage) Store(ctx context.Context, entry *domain.ClipboardEntry) (*domain.ClipboardEntry, error) {
	m.StoreCalled = true
	m.StoreCalledWith = entry
	return m.StoreResult, m.StoreError
}

func (m *MockStorage) GetRecent(ctx context.Context, n int) ([]*domain.ClipboardEntry, error) {
	m.GetRecentCalled = true
	m.GetRecentLimit = n
	return m.GetRecentResult, m.GetRecentError
}

func (m *MockStorage) GetById(ctx context.Context, id string) (*domain.ClipboardEntry, error) {
	m.GetByIdCalled = true
	m.GetByIdId = id
	return m.GetByIdResult, m.GetByIdError
}

func (m *MockStorage) Delete(ctx context.Context, id string) error {
	m.DeleteCalled = true
	m.DeleteId = id
	return m.DeleteError
}

func (m *MockStorage) Search(ctx context.Context, query string, limit int) ([]*domain.ClipboardEntry, error) {
	m.SearchCalled = true
	m.SearchQuery = query
	m.SearchLimit = limit
	return m.SearchResult, m.SearchError
}

func (m *MockStorage) Count(ctx context.Context) (int, error) {
	m.CountCalled = true
	return m.CountResult, m.CountError
}

func (m *MockStorage) Clear(ctx context.Context) error {
	m.ClearCalled = true
	return m.ClearError
}

func (m *MockStorage) DeleteOlderThan(ctx context.Context, cutoff time.Time) (int, error) {
	m.DeleteOlderThanCalled = true
	m.DeleteOlderThanCutoff = cutoff
	return m.DeleteOlderThanResult, m.DeleteOlderThanError
}

type MockAnalyzer struct {
	Result *domain.Analysis
}

func (m *MockAnalyzer) Analyze(entry *domain.ClipboardEntry) *domain.Analysis {
	if m.Result == nil {
		return &domain.Analysis{
			Type:        domain.ContentTypeText,
			IsSensitive: false,
			Reason:      "",
		}
	}
	return m.Result
}

func TestProcessNewEntry(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name            string
		entry           *domain.ClipboardEntry
		analyzerResult  *domain.Analysis
		storageResult   *domain.ClipboardEntry
		storageError    error
		wantErr         error
		wantResult      bool
		wantStoreCalled bool
	}{
		{
			name: "Success",
			entry: &domain.ClipboardEntry{
				Content:   "test content",
				Timestamp: now,
			},
			analyzerResult: &domain.Analysis{
				Type:        domain.ContentTypeText,
				IsSensitive: false,
			},
			storageResult: &domain.ClipboardEntry{
				Id:        "123",
				Content:   "test content",
				Timestamp: now,
			},
			storageError:    nil,
			wantErr:         nil,
			wantResult:      true,
			wantStoreCalled: true,
		},
		{
			name: "SensitiveContent",
			entry: &domain.ClipboardEntry{
				Content:   "password: secret123",
				Timestamp: now,
			},
			analyzerResult: &domain.Analysis{
				Type:        domain.ContentTypeText,
				IsSensitive: true,
				Reason:      "contains password",
			},
			storageResult:   nil,
			storageError:    nil,
			wantErr:         ErrSensitiveContent,
			wantResult:      false,
			wantStoreCalled: false,
		},
		{
			name: "EmptyContent",
			entry: &domain.ClipboardEntry{
				Content:   "",
				Timestamp: now,
			},
			analyzerResult:  nil, // Won't be called
			storageResult:   nil,
			storageError:    nil,
			wantErr:         ErrEmptyContent,
			wantResult:      false,
			wantStoreCalled: false,
		},
		{
			name:            "NilEntry",
			entry:           nil,
			analyzerResult:  nil,
			storageResult:   nil,
			storageError:    nil,
			wantErr:         ErrNilEntry,
			wantResult:      false,
			wantStoreCalled: false,
		},
		{
			name: "StorageError",
			entry: &domain.ClipboardEntry{
				Content:   "valid content",
				Timestamp: now,
			},
			analyzerResult: &domain.Analysis{
				Type:        domain.ContentTypeText,
				IsSensitive: false,
			},
			storageResult:   nil,
			storageError:    errors.New("database error"),
			wantErr:         nil,
			wantResult:      false,
			wantStoreCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &MockStorage{
				StoreResult: tt.storageResult,
				StoreError:  tt.storageError,
			}

			mockAnalyzer := &MockAnalyzer{
				Result: tt.analyzerResult,
			}

			service := NewClipboardService(mockStorage, mockAnalyzer)

			result, err := service.ProcessNewEntry(context.Background(), tt.entry)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.name == "StorageError" {
				if err == nil {
					t.Fatal("expected storage error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}

			if tt.wantResult {
				if result == nil {
					t.Error("expected result, got nil")
				}
			} else {
				if result != nil {
					t.Errorf("expected nil result, got %v", result)
				}
			}

			if mockStorage.StoreCalled != tt.wantStoreCalled {
				t.Errorf("expected StoreCalled=%v, got %v", tt.wantStoreCalled, mockStorage.StoreCalled)
			}

		})
	}
}

func TestGetHistory(t *testing.T) {
	tests := []struct {
		name                string
		limit               int
		storageResult       []*domain.ClipboardEntry
		storageError        error
		wantErr             error
		wantResult          bool
		wantGetRecentCalled bool // Added this
	}{
		{
			name:  "Success",
			limit: 10,
			storageResult: []*domain.ClipboardEntry{
				{Id: "1", Content: "entry1"},
				{Id: "2", Content: "entry2"},
			},
			storageError:        nil,
			wantErr:             nil,
			wantResult:          true,
			wantGetRecentCalled: true, // Should call storage
		},
		{
			name:                "InvalidLimitZero",
			limit:               0,
			storageResult:       nil,
			storageError:        nil,
			wantErr:             ErrInvalidLimit,
			wantResult:          false,
			wantGetRecentCalled: false, // Should NOT call storage
		},
		{
			name:                "InvalidLimitTooHigh",
			limit:               101,
			storageResult:       nil,
			storageError:        nil,
			wantErr:             ErrInvalidLimit,
			wantResult:          false,
			wantGetRecentCalled: false, // Should NOT call storage
		},
		{
			name:                "StorageError",
			limit:               10,
			storageResult:       nil,
			storageError:        errors.New("database error"),
			wantErr:             nil, // Don't check specific error type
			wantResult:          false,
			wantGetRecentCalled: true, // Should call storage (then fail)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &MockStorage{
				GetRecentResult: tt.storageResult,
				GetRecentError:  tt.storageError,
			}

			service := NewClipboardService(mockStorage, &MockAnalyzer{})

			result, err := service.GetHistory(context.Background(), tt.limit)

			// Check error
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.name == "StorageError" {
				if err == nil {
					t.Fatal("expected storage error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}

			if tt.wantResult {
				if len(result) == 0 {
					t.Error("expected result, got nil or empty")
				}
			} else {
				if len(result) > 0 {
					t.Errorf("expected nil or empty result, got %v", result)
				}
			}

			if mockStorage.GetRecentCalled != tt.wantGetRecentCalled {
				t.Errorf("expected GetRecentCalled=%v, got %v",
					tt.wantGetRecentCalled, mockStorage.GetRecentCalled)
			}

			if tt.wantGetRecentCalled && mockStorage.GetRecentLimit != tt.limit {
				t.Errorf("expected GetRecentLimit=%d, got %d",
					tt.limit, mockStorage.GetRecentLimit)
			}
		})
	}
}

func TestGetEntry(t *testing.T) {
	tests := []struct {
		name              string
		id                string
		storageResult     *domain.ClipboardEntry
		storageError      error
		wantErr           error
		wantResult        bool
		wantGetByIdCalled bool
	}{
		{
			name: "Success",
			id:   "123",
			storageResult: &domain.ClipboardEntry{
				Id:      "123",
				Content: "test content",
			},
			storageError:      nil,
			wantErr:           nil,
			wantResult:        true,
			wantGetByIdCalled: true,
		},
		{
			name:              "InvalidId",
			id:                "",
			storageResult:     nil,
			storageError:      nil,
			wantErr:           ErrInvalidId,
			wantResult:        false,
			wantGetByIdCalled: false,
		},
		{
			name:              "NotFound",
			id:                "nonexistent",
			storageResult:     nil,
			storageError:      errors.New("not found"),
			wantErr:           ErrNotFound,
			wantResult:        false,
			wantGetByIdCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &MockStorage{
				GetByIdResult: tt.storageResult,
				GetByIdError:  tt.storageError,
			}

			service := NewClipboardService(mockStorage, &MockAnalyzer{})

			result, err := service.GetEntry(context.Background(), tt.id)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}

			if tt.wantResult {
				if result == nil {
					t.Error("expected result, got nil")
				}
			} else {
				if result != nil {
					t.Errorf("expected nil result, got %v", result)
				}
			}

			if mockStorage.GetByIdCalled != tt.wantGetByIdCalled {
				t.Errorf("expected GetByIdCalled=%v, got %v", tt.wantGetByIdCalled, mockStorage.GetByIdCalled)
			}

			if tt.wantGetByIdCalled && mockStorage.GetByIdId != tt.id {
				t.Errorf("expected GetByIdId=%s, got %s", tt.id, mockStorage.GetByIdId)
			}
		})
	}
}

func TestDeleteEntry(t *testing.T) {
	tests := []struct {
		name             string
		id               string
		storageError     error
		wantErr          error
		wantDeleteCalled bool
	}{
		{
			name:             "Success",
			id:               "123",
			storageError:     nil,
			wantErr:          nil,
			wantDeleteCalled: true,
		},
		{
			name:             "InvalidId",
			id:               "",
			storageError:     nil,
			wantErr:          ErrInvalidId,
			wantDeleteCalled: false,
		},
		{
			name:             "NotFound",
			id:               "nonexistent",
			storageError:     errors.New("not found"),
			wantErr:          ErrNotFound,
			wantDeleteCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &MockStorage{
				DeleteError: tt.storageError,
			}

			service := NewClipboardService(mockStorage, &MockAnalyzer{})

			err := service.DeleteEntry(context.Background(), tt.id)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}

			if mockStorage.DeleteCalled != tt.wantDeleteCalled {
				t.Errorf("expected DeleteCalled=%v, got %v", tt.wantDeleteCalled, mockStorage.DeleteCalled)
			}

			if tt.wantDeleteCalled && mockStorage.DeleteId != tt.id {
				t.Errorf("expected DeleteId=%s, got %s", tt.id, mockStorage.DeleteId)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name             string
		query            string
		limit            int
		storageResult    []*domain.ClipboardEntry
		storageError     error
		wantErr          error
		wantResult       bool
		wantSearchCalled bool
	}{
		{
			name:  "Success",
			query: "test",
			limit: 10,
			storageResult: []*domain.ClipboardEntry{
				{Id: "1", Content: "test content"},
			},
			storageError:     nil,
			wantErr:          nil,
			wantResult:       true,
			wantSearchCalled: true,
		},
		{
			name:             "EmptyQuery",
			query:            "",
			limit:            10,
			storageResult:    nil,
			storageError:     nil,
			wantErr:          ErrEmptyQuery,
			wantResult:       false,
			wantSearchCalled: false,
		},
		{
			name:             "StorageError",
			query:            "test",
			limit:            10,
			storageResult:    nil,
			storageError:     errors.New("database error"),
			wantErr:          nil,
			wantResult:       false,
			wantSearchCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &MockStorage{
				SearchResult: tt.storageResult,
				SearchError:  tt.storageError,
			}

			service := NewClipboardService(mockStorage, &MockAnalyzer{})

			result, err := service.Search(context.Background(), tt.query, tt.limit)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.name == "StorageError" {
				if err == nil {
					t.Fatal("expected storage error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}

			if tt.wantResult {
				if len(result) == 0 {
					t.Error("expected result, got nil or empty")
				}
			} else {
				if len(result) > 0 {
					t.Errorf("expected nil or empty result, got %v", result)
				}
			}

			if mockStorage.SearchCalled != tt.wantSearchCalled {
				t.Errorf("expected SearchCalled=%v, got %v", tt.wantSearchCalled, mockStorage.SearchCalled)
			}

			if tt.wantSearchCalled && mockStorage.SearchQuery != tt.query {
				t.Errorf("expected SearchQuery=%s, got %s", tt.query, mockStorage.SearchQuery)
			}

			if tt.wantSearchCalled && mockStorage.SearchLimit != tt.limit {
				t.Errorf("expected SearchLimit=%d, got %d", tt.limit, mockStorage.SearchLimit)
			}
		})
	}
}

func TestClearHistory(t *testing.T) {
	tests := []struct {
		name            string
		storageError    error
		wantErr         error
		wantClearCalled bool
	}{
		{
			name:            "Success",
			storageError:    nil,
			wantErr:         nil,
			wantClearCalled: true,
		},
		{
			name:            "StorageError",
			storageError:    errors.New("database error"),
			wantErr:         nil,
			wantClearCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &MockStorage{
				ClearError: tt.storageError,
			}

			service := NewClipboardService(mockStorage, &MockAnalyzer{})

			err := service.ClearHistory(context.Background())

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.name == "StorageError" {
				if err == nil {
					t.Fatal("expected storage error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}

			if mockStorage.ClearCalled != tt.wantClearCalled {
				t.Errorf("expected ClearCalled=%v, got %v", tt.wantClearCalled, mockStorage.ClearCalled)
			}
		})
	}
}

func TestGetStats(t *testing.T) {
	tests := []struct {
		name            string
		countResult     int
		storageError    error
		wantErr         error
		wantResult      bool
		wantCountCalled bool
	}{
		{
			name:            "Success",
			countResult:     42,
			storageError:    nil,
			wantErr:         nil,
			wantResult:      true,
			wantCountCalled: true,
		},
		{
			name:            "StorageError",
			countResult:     0,
			storageError:    errors.New("database error"),
			wantErr:         nil,
			wantResult:      false,
			wantCountCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &MockStorage{
				CountResult: tt.countResult,
				CountError:  tt.storageError,
			}

			service := NewClipboardService(mockStorage, &MockAnalyzer{})

			result, err := service.GetStats(context.Background())

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.name == "StorageError" {
				if err == nil {
					t.Fatal("expected storage error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}

			if tt.wantResult {
				if result == nil {
					t.Error("expected result, got nil")
				} else if result.TotalEntries != tt.countResult {
					t.Errorf("expected TotalEntries=%d, got %d", tt.countResult, result.TotalEntries)
				}
			} else {
				if result != nil {
					t.Errorf("expected nil result, got %v", result)
				}
			}

			if mockStorage.CountCalled != tt.wantCountCalled {
				t.Errorf("expected CountCalled=%v, got %v", tt.wantCountCalled, mockStorage.CountCalled)
			}
		})
	}
}

func TestDeleteOlderThan(t *testing.T) {
	now := time.Now()
	cutoff := now.Add(-24 * time.Hour)

	tests := []struct {
		name                      string
		cutoff                    time.Time
		deleteOlderThanResult     int
		deleteOlderThanError      error
		wantErr                   error
		wantResult                int
		wantDeleteOlderThanCalled bool
	}{
		{
			name:                      "Success",
			cutoff:                    cutoff,
			deleteOlderThanResult:     5,
			deleteOlderThanError:      nil,
			wantErr:                   nil,
			wantResult:                5,
			wantDeleteOlderThanCalled: true,
		},
		{
			name:                      "StorageError",
			cutoff:                    cutoff,
			deleteOlderThanResult:     0,
			deleteOlderThanError:      errors.New("database error"),
			wantErr:                   nil,
			wantResult:                0,
			wantDeleteOlderThanCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &MockStorage{
				DeleteOlderThanResult: tt.deleteOlderThanResult,
				DeleteOlderThanError:  tt.deleteOlderThanError,
			}

			service := NewClipboardService(mockStorage, &MockAnalyzer{})

			result, err := service.DeleteOlderThan(context.Background(), tt.cutoff)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else if tt.name == "StorageError" {
				if err == nil {
					t.Fatal("expected storage error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}

			if result != tt.wantResult {
				t.Errorf("expected result %d, got %d", tt.wantResult, result)
			}

			if mockStorage.DeleteOlderThanCalled != tt.wantDeleteOlderThanCalled {
				t.Errorf("expected DeleteOlderThanCalled=%v, got %v", tt.wantDeleteOlderThanCalled, mockStorage.DeleteOlderThanCalled)
			}

			if tt.wantDeleteOlderThanCalled && !mockStorage.DeleteOlderThanCutoff.Equal(tt.cutoff) {
				t.Errorf("expected DeleteOlderThanCutoff=%v, got %v", tt.cutoff, mockStorage.DeleteOlderThanCutoff)
			}
		})
	}
}
