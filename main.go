package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/geodask/clipboard-manager/internal/daemon"
	"github.com/geodask/clipboard-manager/internal/monitor"
	"github.com/geodask/clipboard-manager/internal/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	storage, err := storage.NewSQLiteStorage("./clipboard.db")
	if err != nil {
		fmt.Printf("Failed to initialize storage: %v\n", err)
		return
	}
	defer storage.Close()

	monitor := monitor.NewPollingMonitor()
	daemon := daemon.NewDaemon(monitor, storage)

	errChan := make(chan error, 1)
	go func() {
		errChan <- daemon.Run(ctx)
	}()

	select {
	case <-sigChan:
		fmt.Println("Received shutdown signal")
		cancel()
		<-errChan
		fmt.Println("Shutdown complete")
	case err := <-errChan:
		if err != nil {
			fmt.Printf("Daemon error: %v\n", err)
		}
	}

}
