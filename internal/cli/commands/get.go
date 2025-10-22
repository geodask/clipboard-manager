package commands

import (
	"context"
	"fmt"

	"github.com/geodask/clipboard-manager/internal/client"
)

type GetCommand struct{}

func (c *GetCommand) Name() string {
	return "get"
}

func (c *GetCommand) Description() string {
	return "Get clipboard entry by ID"
}

func (c *GetCommand) Usage() string {
	return "get <id>"
}

func (c *GetCommand) Execute(ctx context.Context, client *client.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Missing required argument: \033[1mid\033[0m\n\n\033[1mUsage:\033[0m\n  \033[2m$\033[0m clipctl \033[36m%s\033[0m\n\n\033[1mExample:\033[0m\n  \033[2m$\033[0m clipctl get abc123\n\n\033[2mTip: Use 'clipctl list' to see available entry IDs\033[0m", c.Usage())
	}

	id := args[0]
	entry, err := client.GetEntry(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving entry: %w", err)
	}

	fmt.Printf("\033[1m┌─ Entry Details\033[0m\n")
	fmt.Printf("\033[1m│\033[0m \033[36mID:\033[0m         %s\n", entry.Id)
	fmt.Printf("\033[1m│\033[0m \033[36mTimestamp:\033[0m  %s\n", entry.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("\033[1m└─ Content:\033[0m\n")
	fmt.Printf("\n%s\n", entry.Content)

	return nil
}
