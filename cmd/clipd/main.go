package main

import (
	"fmt"
	"os"

	"github.com/geodask/clipboard-manager/internal/analyzer"
	"github.com/geodask/clipboard-manager/internal/api"
	"github.com/geodask/clipboard-manager/internal/config"
	"github.com/geodask/clipboard-manager/internal/daemon"
	"github.com/geodask/clipboard-manager/internal/logger"
	"github.com/geodask/clipboard-manager/internal/monitor"
	"github.com/geodask/clipboard-manager/internal/service"
	"github.com/geodask/clipboard-manager/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	logger, err := logger.New(cfg.Logging)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logger.Info("starting clipboard manager daemon", "db_path", cfg.Database.Path, "socket_path", cfg.API.SocketPath, "poll_interval", cfg.Daemon.PollInterval)

	storage, err := storage.NewSQLiteStorage(cfg.Database.Path)
	if err != nil {
		logger.Error("failed to initialize storage", "error", err)
		return
	}
	defer storage.Close()

	monitor := monitor.NewPollingMonitor()
	analyzer := analyzer.NewSimpleAnalyzer()

	service := service.NewClipboardService(storage, analyzer)

	apiServer := api.NewServer(service, cfg.API, logger)

	daemon := daemon.NewDaemon(monitor, service, apiServer, logger, cfg.Daemon)

	if err := daemon.Start(); err != nil {
		logger.Error("daemon stopped with error", "error", err)
		os.Exit(1)
	}

	logger.Info("daemon stopped gracefully")
}
