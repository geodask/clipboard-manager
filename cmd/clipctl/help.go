package main

import (
	"fmt"

	"github.com/geodask/clipboard-manager/internal/cli"
)

func printUsage(registry *cli.Registry) {
	fmt.Printf("%sCLIPCTL - Clipboard Manager CLI%s\n", colorBold, colorReset)
	fmt.Println(separator)
	fmt.Println()
	fmt.Printf("%sUSAGE:%s\n", colorBold, colorReset)
	fmt.Printf("  clipctl [options] %s<command>%s [arguments]\n", colorCyan, colorReset)
	fmt.Println()
	fmt.Printf("%sOPTIONS:%s\n", colorBold, colorReset)
	fmt.Printf("  %s--socket%s <path>     Path to the clipd socket (default: %s)\n", colorCyan, colorReset, defaultSocketPath)
	fmt.Printf("  %s--timeout%s <duration> Request timeout duration (default: %s)\n", colorCyan, colorReset, defaultTimeout)
	fmt.Printf("  %s--version%s           Show version information\n", colorCyan, colorReset)
	fmt.Printf("  %s-v%s                  Verbose output\n", colorCyan, colorReset)
	fmt.Println()
	fmt.Printf("%sCOMMANDS:%s\n", colorBold, colorReset)
	for _, cmd := range registry.All() {
		fmt.Printf("  %s%-12s%s %s\n", colorCyan, cmd.Name(), colorReset, cmd.Description())
		fmt.Printf("               %s$ clipctl %s%s\n", colorDim, cmd.Usage(), colorReset)
		fmt.Println()
	}
	fmt.Printf("  %s%-12s%s %s\n", colorCyan, "help", colorReset, "Show this help message")
	fmt.Printf("               %s$ clipctl help [command]%s\n", colorDim, colorReset)
	fmt.Println()
	fmt.Printf("%sFor more information on a specific command, use: clipctl help <command>%s\n", colorDim, colorReset)
}

func printCommandHelp(cmd cli.Command) {
	fmt.Printf("%s%s%s - %s\n\n", colorBold, cmd.Name(), colorReset, cmd.Description())
	fmt.Printf("%sUSAGE:%s\n", colorBold, colorReset)
	fmt.Printf("  clipctl %s%s%s\n", colorCyan, cmd.Usage(), colorReset)
}

func suggestSimilarCommands(input string, registry *cli.Registry) {
	if len(input) == 0 {
		return
	}

	suggestions := []string{}
	for _, cmd := range registry.All() {
		if len(cmd.Name()) > 0 && input[0] == cmd.Name()[0] {
			suggestions = append(suggestions, cmd.Name())
		}
	}

	if len(suggestions) > 0 {
		fmt.Printf("%sDid you mean one of these?%s\n", colorYellow, colorReset)
		for _, s := range suggestions {
			fmt.Printf("  - %s\n", s)
		}
		fmt.Println()
	}
}
