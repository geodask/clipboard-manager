package daemon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/geodask/clipboard-manager/internal/domain"
	"github.com/geodask/clipboard-manager/internal/service"
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

type Service interface {
	ProcessNewEntry(ctx context.Context, entry *domain.ClipboardEntry) (*domain.ClipboardEntry, error)
}

type Daemon struct {
	monitor Monitor
	service Service
}

func NewDaemon(monitor Monitor, service Service) *Daemon {
	return &Daemon{
		monitor: monitor,
		service: service,
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

				storedEntry, err := d.service.ProcessNewEntry(ctx, entry)

				if err != nil {
					var sensitiveErr *service.SensitiveContentError
					if errors.As(err, &sensitiveErr) {
						fmt.Printf("Skipped sensitive content: %s\n", sensitiveErr.Reason)
					} else {
						fmt.Printf("Error processing entry: %v\n", err)
					}
					continue
				}

				fmt.Printf("Stored entry %s\n", storedEntry.Id)
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
