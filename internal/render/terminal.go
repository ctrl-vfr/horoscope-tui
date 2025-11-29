package render

import (
	"os"

	"golang.org/x/term"
)

// TerminalSize holds the dimensions of the terminal
type TerminalSize struct {
	Width  int
	Height int
}

// GetTerminalSize returns the current terminal dimensions
// Falls back to 80x24 if detection fails
func GetTerminalSize() TerminalSize {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return TerminalSize{Width: 80, Height: 24}
	}
	return TerminalSize{Width: width, Height: height}
}

// MinWidth is the minimum terminal width for full layout
const MinWidth = 100

// MinHeight is the minimum terminal height for full layout
const MinHeight = 30

// TerminalCapabilities describes what the terminal supports
type TerminalCapabilities struct {
	KittyGraphics bool
	Sixel         bool
	TrueColor     bool
}

// GetTerminalCapabilities detects terminal capabilities
func GetTerminalCapabilities() TerminalCapabilities {
	caps := TerminalCapabilities{}

	// Check for Kitty graphics protocol support (Kitty, Ghostty, WezTerm)
	term := os.Getenv("TERM")
	termProgram := os.Getenv("TERM_PROGRAM")
	if os.Getenv("KITTY_WINDOW_ID") != "" ||
		term == "xterm-kitty" ||
		term == "xterm-ghostty" ||
		termProgram == "kitty" ||
		termProgram == "ghostty" ||
		termProgram == "WezTerm" {
		caps.KittyGraphics = true
	}

	// Check for true color support
	colorTerm := os.Getenv("COLORTERM")
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		caps.TrueColor = true
	}

	return caps
}

// CanDisplayImages returns true if the terminal can display images
func CanDisplayImages() bool {
	caps := GetTerminalCapabilities()
	return caps.KittyGraphics && HasResvg()
}
