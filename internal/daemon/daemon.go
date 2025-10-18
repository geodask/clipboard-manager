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
	"golang.org/x/sync/errgroup"
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

type APIServer interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type Daemon struct {
	monitor   Monitor
	service   Service
	apiServer APIServer
}

func NewDaemon(monitor Monitor, service Service, apiServer APIServer) *Daemon {
	return &Daemon{
		monitor:   monitor,
		service:   service,
		apiServer: apiServer,
	}
}

func (d *Daemon) runMonitorLoop(ctx context.Context) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	fmt.Println("Clipboard monitor started")
	for {
		select {
		case <-ticker.C:
			entry, changed, err := d.monitor.Check()
			if err != nil {
				fmt.Printf("Monitor error: %v\n", err)
				continue
			}

			if changed {
				stored, err := d.service.ProcessNewEntry(ctx, entry)
				if err != nil {
					var sensitiveErr *service.SensitiveContentError
					if errors.As(err, &sensitiveErr) {
						fmt.Printf("Skipped sensitive content: %s\n", sensitiveErr.Reason)
					} else {
						fmt.Printf("Error processing entry: %v\n", err)
					}
					continue
				}

				fmt.Printf("Stored entry %s\n", stored.Id)
			}

		case <-ctx.Done():
			fmt.Println("Monitor loop stopping...")
			return ctx.Err()
		}
	}
}

func (d *Daemon) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return d.runMonitorLoop(ctx)
	})

	g.Go(func() error {
		return d.apiServer.Start(ctx)
	})

	g.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return d.apiServer.Shutdown(shutdownCtx)
	})

	return g.Wait()
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
