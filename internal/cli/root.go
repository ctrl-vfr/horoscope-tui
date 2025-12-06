// Package cli contains the root command
package cli

import (
	"github.com/spf13/cobra"

	"github.com/ctrl-vfr/horoscope-tui/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:   "horoscope",
	Short: "Interactive astrological chart TUI",
	Long:  `Interactive terminal application for calculating and visualizing natal charts.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return tui.Run()
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
