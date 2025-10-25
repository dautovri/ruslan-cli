package main

import (
	"fmt"
	"os"

	"github.com/dautovri/ruslan-cli/cmd"
)

var (
	// Set via ldflags during build
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set version info
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date

	// Execute root command
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
