// Package main is the entry point for the horoscope-tui application.
package main

import (
	"log"
	"os"

	"github.com/ctrl-vfr/horoscope-tui/internal/cli"
	"github.com/ctrl-vfr/horoscope-tui/internal/i18n"
)

func init() {
	i18n.Init()

	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}

	if os.Getenv("HOROSCOPE_CITY") == "" {
		log.Fatal("HOROSCOPE_CITY is not set")
	}

	term := os.Getenv("TERM")
	if term != "xterm-kitty" && term != "xterm-ghostty" && term != "wezterm" {
		log.Fatal("Terminal does not support Kitty graphics protocol")
	}
}

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
