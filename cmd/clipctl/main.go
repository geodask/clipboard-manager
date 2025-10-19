package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/geodask/clipboard-manager/internal/client"
)

func main() {
	socketPath := flag.String("socket", "/tmp/clipd.sock", "Path to the clipd socket")
	flag.Parse()

	if flag.NArg() < 1 {
		printUsage()
		return
	}

	command := flag.Arg(0)

	switch command {
	case "list":
		listHistory(*socketPath)
	case "search":
		if flag.NArg() < 2 {
			fmt.Println("Usage: clipctl search <query>")
			return
		}
		searchHistory(*socketPath, flag.Arg(1))
	case "get":
		if flag.NArg() < 2 {
			fmt.Println("Usage: clipctl get <id>")
			return
		}
		getEntry(*socketPath, flag.Arg(1))
	case "delete":
		if flag.NArg() < 1 {
			fmt.Println("Usage: clipctl delete <id>")
			return
		}
		deleteEntry(*socketPath, flag.Arg(1))
	case "stats":
		showStats(*socketPath)
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  clipctl list [n]       - Show last n entries (default 10)")
	fmt.Println("  clipctl search <query> - Search history for query")
	fmt.Println("  clipctl get <id>       - Get specific entry by ID")
	fmt.Println("  clipctl delete <id>    - Delete entry by ID")
	fmt.Println("  clipctl stats          - Show daemon statistics")
}

func listHistory(socketPath string) {
	n := 10
	if flag.NArg() >= 2 {
		if num, err := strconv.Atoi(flag.Arg(1)); err == nil {
			n = num
		}
	}

	// Create client
	c := client.NewClient(socketPath)

	// Check if daemon is running
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := c.Ping(ctx); err != nil {
		fmt.Printf("Error: Daemon not running. Start it with: clipd\n")
		fmt.Printf("Details: %v\n", err)
		return
	}

	// Get history
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	entries, err := c.GetHistory(ctx, n)
	if err != nil {
		fmt.Printf("Error retrieving history: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("No clipboard history found")
		return
	}

	fmt.Printf("Last %d clipboard entries:\n\n", len(entries))

	// Display in reverse order (most recent last, like before)
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		fmt.Printf("[%s] (ID: %s)\n%s\n---\n",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Id,
			truncate(entry.Content, 100))
	}
}

func searchHistory(socketPath string, query string) {
	// Create client
	c := client.NewClient(socketPath)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	entries, err := c.Search(ctx, query, 50)
	if err != nil {
		fmt.Printf("Error searching: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Printf("No entries found matching '%s'\n", query)
		return
	}

	fmt.Printf("Found %d entries matching '%s':\n\n", len(entries), query)

	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		fmt.Printf("[%s] (ID: %s)\n%s\n---\n",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Id,
			truncate(entry.Content, 100))
	}
}

func getEntry(socketPath string, id string) {
	// Create client
	c := client.NewClient(socketPath)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	entry, err := c.GetEntry(ctx, id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entry ID: %s\n", entry.Id)
	fmt.Printf("Timestamp: %s\n", entry.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("Content:\n%s\n", entry.Content)
}

func deleteEntry(socketPath string, id string) {
	// Create client
	c := client.NewClient(socketPath)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.DeleteEntry(ctx, id); err != nil {
		fmt.Printf("Error deleting entry: %v\n", err)
		return
	}

	fmt.Printf("Entry %s deleted successfully\n", id)
}

func showStats(socketPath string) {
	// Create client
	c := client.NewClient(socketPath)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats, err := c.GetStats(ctx)
	if err != nil {
		fmt.Printf("Error getting stats: %v\n", err)
		return
	}

	fmt.Println("Daemon Statistics:")
	fmt.Printf("  Status: %s\n", stats.Status)
	fmt.Printf("  Total Entries: %d\n", stats.TotalEntries)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
