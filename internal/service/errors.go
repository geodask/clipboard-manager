package service

import "errors"

var (
	// Entry-related errors
	ErrNotFound     = errors.New("entry not found")
	ErrInvalidId    = errors.New("invalid entry ID")
	ErrEmptyContent = errors.New("content cannot empty")
	ErrNilEntry     = errors.New("entry cannot be nil")

	// Query-related errors
	ErrInvalidLimit = errors.New("limit must be between 1 and 1000")
	ErrEmptyQuery   = errors.New("search query cannot be empty")

	// Content-related errors
	ErrSensitiveContent = errors.New("content contains sensitive data")
)

type SensitiveContentError struct {
	Reason string
}

func (e *SensitiveContentError) Error() string {
	return "sensitive content: " + e.Reason
}

func (e *SensitiveContentError) Is(target error) bool {
	return target == ErrSensitiveContent
}
