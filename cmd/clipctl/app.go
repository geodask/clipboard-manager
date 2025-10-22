package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/geodask/clipboard-manager/internal/cli"
	"github.com/geodask/clipboard-manager/internal/client"
)

const (
	exitSuccess          = 0
	exitGeneralError     = 1
	exitCommandNotFound  = 2
	exitDaemonNotRunning = 3
)

type app struct {
	config   *config
	registry *cli.Registry
	output   *output
}

func (a *app) run() int {
	if a.config.showVersion {
		a.output.printVersion()
		return exitSuccess
	}

	if flag.NArg() < 1 {
		printUsage(a.registry)
		return exitSuccess
	}

	commandName := flag.Arg(0)

	if commandName == "help" {
		return a.handleHelpCommand()
	}

	cmd, ok := a.registry.Get(commandName)
	if !ok {
		return a.handleUnknownCommand(commandName)
	}

	return a.executeCommand(cmd, commandName)
}

// handleHelpCommand handles the 'help' command
func (a *app) handleHelpCommand() int {
	if flag.NArg() > 1 {
		cmdName := flag.Arg(1)
		if cmd, ok := a.registry.Get(cmdName); ok {
			printCommandHelp(cmd)
			return exitSuccess
		}
		a.output.printError(fmt.Sprintf("Unknown command: %s", cmdName))
		fmt.Println()
		printUsage(a.registry)
		return exitCommandNotFound
	}
	printUsage(a.registry)
	return exitSuccess
}

// handleUnknownCommand handles unknown command errors
func (a *app) handleUnknownCommand(commandName string) int {
	a.output.printError(fmt.Sprintf("Unknown command: %s", commandName))
	fmt.Println()
	suggestSimilarCommands(commandName, a.registry)
	printUsage(a.registry)
	return exitCommandNotFound
}

// executeCommand executes the given command with proper setup
func (a *app) executeCommand(cmd cli.Command, commandName string) int {
	ctx := a.setupSignalHandling()

	client := client.NewClient(a.config.socketPath)

	if err := a.pingDaemon(ctx, client); err != nil {
		return exitDaemonNotRunning
	}

	a.output.logVerbose("Executing: %s", commandName)

	cmdCtx, cancel := context.WithTimeout(ctx, a.config.timeout)
	defer cancel()

	if err := cmd.Execute(cmdCtx, client, flag.Args()[1:]); err != nil {
		a.output.printError(err.Error())
		return exitGeneralError
	}

	a.output.logSuccess("Command completed successfully")
	return exitSuccess
}

// setupSignalHandling sets up graceful shutdown on interrupt signals
func (a *app) setupSignalHandling() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		a.output.logVerbose("\nReceived interrupt signal, shutting down...")
		cancel()
		os.Exit(exitGeneralError)
	}()

	return ctx
}

// pingDaemon checks if the daemon is running
func (a *app) pingDaemon(ctx context.Context, c *client.Client) error {
	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := c.Ping(pingCtx); err != nil {
		a.output.printDaemonError(err, a.config.socketPath)
		return err
	}
	return nil
}
