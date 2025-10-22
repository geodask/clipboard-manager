package main

import (
	"github.com/geodask/clipboard-manager/internal/cli"
	"github.com/geodask/clipboard-manager/internal/cli/commands"
)

func createRegistry() *cli.Registry {
	registry := cli.NewRegistry()
	registry.Register(&commands.ListCommand{})
	registry.Register(&commands.SearchCommand{})
	registry.Register(&commands.GetCommand{})
	registry.Register(&commands.DeleteCommand{})
	registry.Register(&commands.StatsCommand{})
	return registry
}
