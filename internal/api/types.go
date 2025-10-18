package api

import "time"

type GetHistoryRequest struct {
	Limit int `json:"limit"`
}

type SearchRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

type CreateEntryRequest struct {
	Content string `json:"content"`
}

type EntryResponse struct {
	Id        string    `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type HistoryResponse struct {
	Entries []EntryResponse `json:"entries"`
	Total   int             `json:"total"`
}

type StatsResponse struct {
	TotalEntries int    `json:"total_entries"`
	Status       string `json:"status"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
