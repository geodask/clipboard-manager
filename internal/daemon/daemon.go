package daemon

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/geodask/clipboard-manager/internal/domain"
)

type Monitor interface {
	Check() (entry *domain.ClipboardEntry, changed bool, err error)
}

type Storage interface {
	Store(ctx context.Context, entry *domain.ClipboardEntry) (*domain.ClipboardEntry, error)
}

type Analyzer interface {
	Analyze(entry *domain.ClipboardEntry) *domain.Analysis
}

type Daemon struct {
	monitor  Monitor
	storage  Storage
	analyzer Analyzer
}

func NewDaemon(monitor Monitor, storage Storage, analyzer Analyzer) *Daemon {
	return &Daemon{
		monitor:  monitor,
		storage:  storage,
		analyzer: analyzer,
	}
}

func (d *Daemon) Run(ctx context.Context) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			entry, changed, err := d.monitor.Check()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			if changed {
				analysis := d.analyzer.Analyze(entry)

				if analysis.IsSensitive {
					fmt.Printf("Sensitive content detected: %s\n", analysis.Reason)
					continue
				}

				fmt.Printf("Storing %s content\n", analysis.Type)

				storedEntry, err := d.storage.Store(ctx, entry)
				if err != nil {
					if ctx.Err() != nil {
						fmt.Println("Storage cancelled - shutting down")
						return ctx.Err()
					}
					fmt.Printf("Storage error: %v\n", err)
				} else {
					fmt.Printf("Stored with ID: %s\n", storedEntry.Id)
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (d *Daemon) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		errChan <- d.Run(ctx)
	}()

	select {
	case <-sigChan:
		fmt.Println("\nShutting down gracefully...")
		cancel()
		err := <-errChan
		fmt.Println("Shutdown complete.")
		return err
	case err := <-errChan:
		return err
	}

}
