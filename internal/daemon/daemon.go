package daemon

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/geodask/clipboard-manager/internal/config"
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
	DeleteOlderThan(ctx context.Context, cutoff time.Time) (int, error)
}

type APIServer interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type Daemon struct {
	monitor           Monitor
	service           Service
	apiServer         APIServer
	pollInterval      time.Duration
	shutdownTimeout   time.Duration
	retentionEnabled  bool
	retentionMaxAge   time.Duration
	retentionInterval time.Duration
	logger            *slog.Logger
}

func NewDaemon(
	monitor Monitor,
	service Service,
	apiServer APIServer,
	logger *slog.Logger,
	cfg config.DaemonConfig,
) *Daemon {
	return &Daemon{
		monitor:           monitor,
		service:           service,
		apiServer:         apiServer,
		pollInterval:      cfg.PollInterval,
		shutdownTimeout:   cfg.ShutdownTimeout,
		retentionEnabled:  cfg.RetentionEnabled,
		retentionMaxAge:   cfg.RetentionMaxAge,
		retentionInterval: cfg.RetentionInterval,
		logger:            logger,
	}
}

func (d *Daemon) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return d.runMonitorLoop(ctx)
	})

	g.Go(func() error {
		return d.runRetentionLoop(ctx)
	})

	g.Go(func() error {
		return d.apiServer.Start(ctx)
	})

	g.Go(func() error {
		<-ctx.Done()

		d.logger.Info("initiating graceful shutdown", "timeout", d.shutdownTimeout)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), d.shutdownTimeout)
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
		d.logger.Info("received shutdown signal")
		cancel()
		err := <-errChan
		d.logger.Info("shutdown complete")
		return err
	case err := <-errChan:
		if err != nil {
			d.logger.Error("daemon error", "error", err)
		}
		return err
	}

}

func (d *Daemon) PerformRetention(ctx context.Context) (int, error) {
	cutoff := time.Now().Add(-d.retentionMaxAge)
	deleted, err := d.service.DeleteOlderThan(ctx, cutoff)
	return deleted, err
}

func (d *Daemon) runMonitorLoop(ctx context.Context) error {
	ticker := time.NewTicker(d.pollInterval)
	defer ticker.Stop()

	d.logger.Info("clipboard monitor started", "poll_interval", d.pollInterval)
	for {
		select {
		case <-ticker.C:
			entry, changed, err := d.monitor.Check()
			if err != nil {
				d.logger.Error("monitor check failed", "error", err)
				continue
			}
			if changed {
				stored, err := d.service.ProcessNewEntry(ctx, entry)
				if err != nil {
					var sensitiveErr *service.SensitiveContentError
					if errors.As(err, &sensitiveErr) {
						d.logger.Debug("skipped sensitive content", "reason", sensitiveErr.Reason, "content_length", len(entry.Content))
					} else {
						d.logger.Error("failed to process entry", "error", err, "content_length", len(entry.Content))
					}
					continue
				}

				d.logger.Info("stored clipboard entry", "id", stored.Id, "content_length", len(stored.Content), "timestamp", stored.Timestamp)
			}

		case <-ctx.Done():
			d.logger.Info("monitor loop stopping")
			return ctx.Err()
		}
	}
}

func (d *Daemon) runRetentionLoop(ctx context.Context) error {
	if !d.retentionEnabled {
		d.logger.Info("retention cleanup disabled")
		<-ctx.Done()
		return ctx.Err()
	}

	d.logger.Info("performing initial retention cleanup")
	deleted, err := d.PerformRetention(ctx)
	if err != nil {
		d.logger.Error("initial retention cleanup failed", "error", err)
	} else {
		d.logger.Info("initial retention cleanup completed", "deleted_entries", deleted)
	}

	ticker := time.NewTicker(d.retentionInterval)
	defer ticker.Stop()

	d.logger.Info("retention cleanup started", "interval", d.retentionInterval)
	for {
		select {
		case <-ticker.C:
			deleted, err := d.PerformRetention(ctx)
			if err != nil {
				d.logger.Error("retention cleanup failed", "error", err)
				continue
			}
			d.logger.Info("retention cleanup completed", "deleted_entries", deleted)
		case <-ctx.Done():
			d.logger.Info("retention loop stopping")
			return ctx.Err()
		}
	}
}
