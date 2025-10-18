package monitor

import (
	"time"

	"github.com/atotto/clipboard"
	"github.com/geodask/clipboard-manager/internal/domain"
)

type Monitor interface {
	Check() (entry *domain.ClipboardEntry, changed bool, err error)
}

type PollingMonitor struct {
	lastContent string
}

func NewPollingMonitor() *PollingMonitor {
	return &PollingMonitor{}
}

func (pm *PollingMonitor) Check() (*domain.ClipboardEntry, bool, error) {
	content, err := clipboard.ReadAll()

	if err != nil {
		return nil, false, err
	}

	changed := content != pm.lastContent && content != ""

	if changed {
		pm.lastContent = content
		entry := &domain.ClipboardEntry{
			Content:   content,
			Timestamp: time.Now(),
		}
		return entry, true, nil
	}

	return nil, false, nil
}
