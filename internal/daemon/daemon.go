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
	Store(entry *domain.ClipboardEntry) error
}

type Daemon struct {
	monitor Monitor
	storage Storage
}

func NewDaemon(monitor Monitor, storage Storage) *Daemon {
	return &Daemon{
		monitor: monitor,
		storage: storage,
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
				if err := d.storage.Store(entry); err != nil {
					fmt.Printf("Storage error: %v\n", err)
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
