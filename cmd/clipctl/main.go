package main

import (
	"os"
)

const version = "1.0.0"

func main() {
	cfg := parseFlags()
	registry := createRegistry()

	app := &app{
		config:   cfg,
		registry: registry,
		output:   newOutput(cfg.verbose),
	}

	os.Exit(app.run())
}
