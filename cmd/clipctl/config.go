package main

import (
	"flag"
	"time"
)

const (
	defaultSocketPath = "/tmp/clipd.sock"
	defaultTimeout    = 5 * time.Second
	pingTimeout       = 2 * time.Second
)

type config struct {
	socketPath  string
	timeout     time.Duration
	showVersion bool
	verbose     bool
}

func parseFlags() *config {
	cfg := &config{}

	flag.StringVar(&cfg.socketPath, "socket", defaultSocketPath, "Path to the clipd socket")
	flag.DurationVar(&cfg.timeout, "timeout", defaultTimeout, "Request timeout duration")
	flag.BoolVar(&cfg.showVersion, "version", false, "Show version information")
	flag.BoolVar(&cfg.verbose, "v", false, "Verbose output")

	flag.Usage = func() {
		registry := createRegistry()
		printUsage(registry)
	}

	flag.Parse()
	return cfg
}
