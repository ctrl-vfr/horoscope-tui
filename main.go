package main

import (
	"log"
	"os"

	"github.com/ctrl-vfr/horoscope-tui/internal/cli"
)

func init() {
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
