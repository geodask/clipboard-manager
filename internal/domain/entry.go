package domain

import "time"

type ClipboardEntry struct {
	Content   string
	Timestamp time.Time
}
