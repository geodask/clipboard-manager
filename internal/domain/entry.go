package domain

import "time"

type ClipboardEntry struct {
	Id        string
	Content   string
	Timestamp time.Time
}

type ContentType string

const (
	ContentTypeText     ContentType = "text"
	ContentTypeURL      ContentType = "url"
	ContentTypeCode     ContentType = "code"
	ContentTypeFilePath ContentType = "filepath"
	ContentTypeUknown   ContentType = "unknown"
)

type Analysis struct {
	Type        ContentType
	IsSensitive bool
	Reason      string
}
