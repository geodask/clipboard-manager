package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/geodask/clipboard-manager/internal/client"
)

type ListCommand struct{}

func (c *ListCommand) Name() string {
	return "list"
}

func (c *ListCommand) Description() string {
	return "Show last n entries (default 10)"
}

func (c *ListCommand) Usage() string {
	return "list [n]"
}

func (c *ListCommand) Execute(ctx context.Context, client *client.Client, args []string) error {
	n := 10
	if len(args) > 0 {
		if num, err := strconv.Atoi(args[0]); err == nil {
			n = num
		}
	}

	entries, err := client.GetHistory(ctx, n)
	if err != nil {
		return fmt.Errorf("retrieving history: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No clipboard history found")
		return nil
	}

	fmt.Printf("\033[1mLast %d clipboard entries:\033[0m\n\n", len(entries))
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		fmt.Printf("\033[2m[\033[0m\033[36m%s\033[0m\033[2m]\033[0m \033[2m(ID: %s)\033[0m\n%s\n\033[2m───────────────────────────────────────────────────────────────\033[0m\n",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Id,
			truncate(entry.Content, 100))
	}

	return nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
