package main

import (
	"fmt"
	"os"
)

type output struct {
	verbose bool
}

func newOutput(verbose bool) *output {
	return &output{verbose: verbose}
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
	colorYellow = "\033[33m"
	separator   = "\033[2m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m"
)

func (o *output) printVersion() {
	fmt.Printf("%sclipctl%s version %s%s%s\n", colorBold, colorReset, colorCyan, version, colorReset)
	fmt.Printf("%sClipboard Manager CLI Tool%s\n", colorDim, colorReset)
}

func (o *output) printError(msg string) {
	fmt.Fprintf(os.Stderr, "%sError:%s %s\n", colorRed, colorReset, msg)
}

func (o *output) printDaemonError(err error, socketPath string) {
	fmt.Printf("%sError:%s Daemon not running\n\n", colorRed, colorReset)
	fmt.Printf("%sTo start the daemon:%s\n", colorBold, colorReset)
	fmt.Printf("  %s$%s clipd\n\n", colorDim, colorReset)

	if o.verbose {
		fmt.Printf("%sDetails: %v%s\n", colorDim, err, colorReset)
		fmt.Printf("%sSocket path: %s%s\n", colorDim, socketPath, colorReset)
	}
}

func (o *output) logVerbose(format string, args ...interface{}) {
	if o.verbose {
		fmt.Printf(colorDim+format+colorReset+"\n", args...)
	}
}

func (o *output) logSuccess(msg string) {
	if o.verbose {
		fmt.Printf("%s✓ %s%s\n", colorDim, msg, colorReset)
	}
}
