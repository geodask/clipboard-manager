package main

import (
	"fmt"

	"github.com/geodask/clipboard-manager/internal/analyzer"
	"github.com/geodask/clipboard-manager/internal/api"
	"github.com/geodask/clipboard-manager/internal/daemon"
	"github.com/geodask/clipboard-manager/internal/monitor"
	"github.com/geodask/clipboard-manager/internal/service"
	"github.com/geodask/clipboard-manager/internal/storage"
)

func main() {
	storage, err := storage.NewSQLiteStorage("./clipboard.db")
	if err != nil {
		fmt.Printf("Failed to initialize storage: %v\n", err)
		return
	}
	defer storage.Close()

	monitor := monitor.NewPollingMonitor()
	analyzer := analyzer.NewSimpleAnalyzer()

	service := service.NewClipboardService(storage, analyzer)

	apiServer := api.NewServer(service, "/tmp/clipd.sock")

	daemon := daemon.NewDaemon(monitor, service, apiServer)

	if err := daemon.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

}
