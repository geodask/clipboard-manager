package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/geodask/clipboard-manager/internal/domain"
	"github.com/geodask/clipboard-manager/internal/service"
	"github.com/go-chi/chi/v5"
)

type Service interface {
	ProcessNewEntry(ctx context.Context, entry *domain.ClipboardEntry) (*domain.ClipboardEntry, error)
	GetHistory(ctx context.Context, limit int) ([]*domain.ClipboardEntry, error)
	GetEntry(ctx context.Context, id string) (*domain.ClipboardEntry, error)
	DeleteEntry(ctx context.Context, id string) error
	Search(ctx context.Context, query string, limit int) ([]*domain.ClipboardEntry, error)
	ClearHistory(ctx context.Context) error
	GetStats(ctx context.Context) (*service.Stats, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	entries, err := h.service.GetHistory(r.Context(), limit)
	if err != nil {
		respondError(w, err)
	}

	var entryResponses []EntryResponse
	for _, entry := range entries {
		entryResponses = append(entryResponses, EntryResponse{
			Id:        entry.Id,
			Content:   entry.Content,
			Timestamp: entry.Timestamp,
		})
	}

	respondJSON(w, http.StatusOK, HistoryResponse{
		Entries: entryResponses,
		Total:   len(entryResponses),
	})
}

func (h *Handler) GetEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		respondError(w, service.ErrInvalidId)
		return
	}

	entry, err := h.service.GetEntry(r.Context(), id)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, EntryResponse{
		Id:        entry.Id,
		Content:   entry.Content,
		Timestamp: entry.Timestamp,
	})
}

func (h *Handler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		respondError(w, service.ErrInvalidId)
		return
	}

	err := h.service.DeleteEntry(r.Context(), id)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{
		Message: "Entry deleted successfully",
	})
}

// POST /api/v1/entries
func (h *Handler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var req CreateEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid JSON",
		})
		return
	}

	// Create entry
	entry := &domain.ClipboardEntry{
		Content:   req.Content,
		Timestamp: time.Now(),
	}

	stored, err := h.service.ProcessNewEntry(r.Context(), entry)
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, EntryResponse{
		Id:        stored.Id,
		Content:   stored.Content,
		Timestamp: stored.Timestamp,
	})
}

// GET /api/v1/search?q=query&limit=10
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")

	limit := 100 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	entries, err := h.service.Search(r.Context(), query, limit)
	if err != nil {
		respondError(w, err)
		return
	}

	var entryResponses []EntryResponse
	for _, entry := range entries {
		entryResponses = append(entryResponses, EntryResponse{
			Id:        entry.Id,
			Content:   entry.Content,
			Timestamp: entry.Timestamp,
		})
	}

	respondJSON(w, http.StatusOK, HistoryResponse{
		Entries: entryResponses,
		Total:   len(entryResponses),
	})
}

// DELETE /api/v1/history
func (h *Handler) ClearHistory(w http.ResponseWriter, r *http.Request) {
	err := h.service.ClearHistory(r.Context())
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{
		Message: "History cleared successfully",
	})
}

// GET /api/v1/stats
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		respondError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, StatsResponse{
		TotalEntries: stats.TotalEntries,
		Status:       "running",
	})
}
