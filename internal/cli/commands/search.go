package commands

import (
	"context"
	"fmt"

	"github.com/geodask/clipboard-manager/internal/client"
)

type SearchCommand struct{}

func (c *SearchCommand) Name() string {
	return "search"
}

func (c *SearchCommand) Description() string {
	return "Search history for query"
}

func (c *SearchCommand) Usage() string {
	return "search <query>"
}

func (c *SearchCommand) Execute(ctx context.Context, client *client.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Missing required argument: \033[1mquery\033[0m\n\n\033[1mUsage:\033[0m\n  \033[2m$\033[0m clipctl \033[36m%s\033[0m\n\n\033[1mExample:\033[0m\n  \033[2m$\033[0m clipctl search \"password\"\n  \033[2m$\033[0m clipctl search code", c.Usage())
	}

	query := args[0]
	entries, err := client.Search(ctx, query, 50)
	if err != nil {
		return fmt.Errorf("searching: %w", err)
	}

	if len(entries) == 0 {
		fmt.Printf("No entries found matching \033[1m'%s'\033[0m\n", query)
		return nil
	}

	fmt.Printf("Found \033[1m%d\033[0m entries matching \033[1m'%s'\033[0m:\n\n", len(entries), query)
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		fmt.Printf("\033[2m[\033[0m\033[36m%s\033[0m\033[2m]\033[0m \033[2m(ID: %s)\033[0m\n%s\n\033[2m───────────────────────────────────────────────────────────────\033[0m\n",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Id,
			truncate(entry.Content, 100))
	}

	return nil
}
