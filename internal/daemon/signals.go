package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type SignalHandler struct {
	daemon  *Daemon
	sigChan chan os.Signal
}

func NewSignalHandler(daemon *Daemon) *SignalHandler {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	return &SignalHandler{
		daemon:  daemon,
		sigChan: sigChan,
	}
}
func (sh *SignalHandler) Handle(ctx context.Context, sig os.Signal) bool {
	switch sig {
	case os.Interrupt, syscall.SIGTERM:
		sh.daemon.logger.Info("received shutdown signal", "signal", sig)
		return false

	case syscall.SIGHUP:
		sh.daemon.logger.Info("received reload signal")
		sh.handleReload()
		return true

	case syscall.SIGUSR1:
		sh.daemon.logger.Info("received manual retention trigger")
		sh.handleManualRetention(ctx)
		return true

	case syscall.SIGUSR2:
		sh.daemon.logger.Info("received stats dump request")
		sh.handleStatsDump(ctx)
		return true
	}

	return true
}

func (sh *SignalHandler) handleReload() {

	sh.daemon.logger.Info("current configuration",
		"poll_interval", sh.daemon.pollInterval,
		"retention_enabled", sh.daemon.retentionEnabled,
		"retention_max_age", sh.daemon.retentionMaxAge,
		"retention_interval", sh.daemon.retentionInterval,
	)

	sh.daemon.logger.Info("note: full config reload requires daemon restart")
}

func (sh *SignalHandler) handleManualRetention(ctx context.Context) {
	deleted, err := sh.daemon.PerformRetention(ctx)
	if err != nil {
		sh.daemon.logger.Error("manual retention failed", "error", err)
		return
	}
	sh.daemon.logger.Info("manual retention completed", "deleted_entries", deleted)

}

func (sh *SignalHandler) handleStatsDump(ctx context.Context) {
	uptime := time.Since(sh.daemon.startTime)

	sh.daemon.logger.Info("daemon statistics",
		"uptime", uptime.Round(time.Second),
		"pid", os.Getpid(),
	)

}

func (sh *SignalHandler) Stop() {
	signal.Stop(sh.sigChan)
	close(sh.sigChan)
}
