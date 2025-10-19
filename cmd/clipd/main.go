package main

import (
	"fmt"
	"os"

	"github.com/geodask/clipboard-manager/internal/analyzer"
	"github.com/geodask/clipboard-manager/internal/api"
	"github.com/geodask/clipboard-manager/internal/config"
	"github.com/geodask/clipboard-manager/internal/daemon"
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

	storage, err := storage.NewSQLiteStorage(cfg.Database.Path)
	if err != nil {
		fmt.Printf("Failed to initialize storage: %v\n", err)
		return
	}
	defer storage.Close()

	monitor := monitor.NewPollingMonitor()
	analyzer := analyzer.NewSimpleAnalyzer()

	service := service.NewClipboardService(storage, analyzer)

	apiServer := api.NewServer(service, cfg.API)

	daemon := daemon.NewDaemon(monitor, service, apiServer, cfg.Daemon.PollInterval, cfg.Daemon.ShutdownTimeout)

	if err := daemon.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

}
