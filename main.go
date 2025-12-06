// Package main is the entry point for the astral-tui application.
package main

import (
	"fmt"
	"os"

	"github.com/ctrl-vfr/astral-tui/internal/cli"
	"github.com/ctrl-vfr/astral-tui/internal/i18n"
	"github.com/ctrl-vfr/astral-tui/internal/preflight"
)

func main() {
	i18n.Init()

	results := preflight.RunChecks()
	if !preflight.PrintResults(results) {
		fmt.Println()
		os.Exit(1)
	}

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
