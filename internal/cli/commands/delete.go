package commands

import (
	"context"
	"fmt"

	"github.com/geodask/clipboard-manager/internal/client"
)

type DeleteCommand struct{}

func (c *DeleteCommand) Name() string {
	return "delete"
}

func (c *DeleteCommand) Description() string {
	return "Delete clipboard entry by ID"
}

func (c *DeleteCommand) Usage() string {
	return "delete <id>"
}

func (c *DeleteCommand) Execute(ctx context.Context, client *client.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Missing required argument: \033[1mid\033[0m\n\n\033[1mUsage:\033[0m\n  \033[2m$\033[0m clipctl \033[36m%s\033[0m\n\n\033[1mExample:\033[0m\n  \033[2m$\033[0m clipctl delete abc123\n\n\033[2mTip: Use 'clipctl list' to see available entry IDs\033[0m", c.Usage())
	}

	id := args[0]
	if err := client.DeleteEntry(ctx, id); err != nil {
		return fmt.Errorf("deleting entry: %w", err)
	}

	fmt.Printf("Entry \033[1m%s\033[0m deleted successfully\n", id)
	return nil
}
