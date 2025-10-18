package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/geodask/clipboard-manager/internal/storage"
)

func main() {

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "list":
		listHistory()
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Usage: clipctl search <query>")
			return
		}
		searchHistory(os.Args[2])
	default:
		printUsage()
	}

}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  clipctl list [n]       - Show last n entries (default 10)")
	fmt.Println("  clipctl search <query> - Search history for query")
}

func listHistory() {
	n := 10
	if len(os.Args) >= 3 {
		if num, err := strconv.Atoi(os.Args[2]); err == nil {
			n = num
		}
	}

	storage, err := storage.NewSQLiteStorage("./clipboard.db")
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
		return
	}
	defer storage.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	entries, err := storage.GetRecent(ctx, n)
	defer cancel()

	if err != nil {
		fmt.Printf("Error retrieving history: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("No clipboard history found")
		return
	}

	fmt.Printf("Last %d clipboard entries:\n\n", len(entries))
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		fmt.Printf("[%s] ID: %s\n%s\n---\n", // Add ID display
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Id,
			truncate(entry.Content, 100))
	}

}

func searchHistory(query string) {
	fmt.Printf("Searching for: %s\n", query)
	fmt.Println("(Search not implemented yet - use: clipctl list | grep \"your query\")")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
