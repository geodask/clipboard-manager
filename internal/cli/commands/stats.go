package commands

import (
	"context"
	"fmt"

	"github.com/geodask/clipboard-manager/internal/client"
)

type StatsCommand struct{}

func (c *StatsCommand) Name() string {
	return "stats"
}

func (c *StatsCommand) Description() string {
	return "Show clipboard statistics"
}

func (c *StatsCommand) Usage() string {
	return "stats"
}

func (c *StatsCommand) Execute(ctx context.Context, client *client.Client, args []string) error {
	stats, err := client.GetStats(ctx)
	if err != nil {
		return fmt.Errorf("Error getting stats: %v\n", err)
	}

	fmt.Println("\033[1m┌─ Daemon Statistics\033[0m")
	fmt.Printf("\033[1m│\033[0m \033[36mStatus:\033[0m         %s\n", stats.Status)
	fmt.Printf("\033[1m└─\033[0m \033[36mTotal Entries:\033[0m %d\n", stats.TotalEntries)

	return nil
}
