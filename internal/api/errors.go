package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/geodask/clipboard-manager/internal/service"
)

func respondError(w http.ResponseWriter, err error) {
	var statusCode int
	var message string

	switch {
	case errors.Is(err, service.ErrNotFound):
		statusCode = http.StatusNotFound
		message = "Entry not found"

	case errors.Is(err, service.ErrInvalidId):
		statusCode = http.StatusBadRequest
		message = "Invalid entry ID"

	case errors.Is(err, service.ErrInvalidLimit):
		statusCode = http.StatusBadRequest
		message = "Invalid limit parameter"

	case errors.Is(err, service.ErrEmptyQuery):
		statusCode = http.StatusBadRequest
		message = "Search query cannot be empty"

	case errors.Is(err, service.ErrEmptyContent):
		statusCode = http.StatusBadRequest
		message = "Content cannot be empty"

	case errors.Is(err, service.ErrSensitiveContent):
		statusCode = http.StatusBadRequest
		message = "Content contains sensitive data"

	default:
		statusCode = http.StatusInternalServerError
		message = "Internal server error"
	}

	respondJSON(w, statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
