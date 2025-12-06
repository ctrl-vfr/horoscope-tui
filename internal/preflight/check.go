// Package preflight provides dependency checking before application startup.
package preflight

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/lipgloss"
	"github.com/ctrl-vfr/astral-tui/internal/i18n"
)

// CheckResult represents the result of a single dependency check.
type CheckResult struct {
	Name string
	OK   bool
	Help string
}

// RunChecks verifies all required dependencies and returns the results.
func RunChecks() []CheckResult {
	return []CheckResult{
		checkOpenAIKey(),
		checkCity(),
		checkTerminal(),
		checkResvg(),
	}
}

func checkOpenAIKey() CheckResult {
	ok := os.Getenv("OPENAI_API_KEY") != ""
	return CheckResult{
		Name: "OPENAI_API_KEY",
		OK:   ok,
		Help: i18n.T("PreflightOpenAIHelp"),
	}
}

func checkCity() CheckResult {
	ok := os.Getenv("ASTRAL_CITY") != ""
	return CheckResult{
		Name: "ASTRAL_CITY",
		OK:   ok,
		Help: i18n.T("PreflightCityHelp"),
	}
}

func checkTerminal() CheckResult {
	term := os.Getenv("TERM")
	ok := term == "xterm-kitty" || term == "xterm-ghostty" || term == "wezterm"
	return CheckResult{
		Name: i18n.T("PreflightTerminal"),
		OK:   ok,
		Help: i18n.T("PreflightTerminalHelp"),
	}
}

func checkResvg() CheckResult {
	_, err := exec.LookPath("resvg")
	return CheckResult{
		Name: "resvg",
		OK:   err == nil,
		Help: i18n.T("PreflightResvgHelp"),
	}
}

// PrintResults displays check results only if there are failures.
func PrintResults(results []CheckResult) bool {
	allOK := true
	for _, r := range results {
		if !r.OK {
			allOK = false
			break
		}
	}

	if allOK {
		return true
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	for _, r := range results {
		if r.OK {
			fmt.Println(successStyle.Render("✓") + " " + r.Name)
		} else {
			fmt.Println(errorStyle.Render("✗") + " " + r.Name)
			fmt.Println(helpStyle.Render("  → " + r.Help))
			fmt.Println()
		}
	}

	return false
}
