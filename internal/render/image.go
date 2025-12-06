// Package render provides SVG to PNG conversion and terminal graphics support.
package render

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// SVGToPNG converts SVG data to PNG using resvg CLI
func SVGToPNG(svgData []byte, width, height int) ([]byte, error) {
	// Check if resvg is available
	resvgPath, err := exec.LookPath("resvg")
	if err != nil {
		return nil, fmt.Errorf("resvg not found: install with 'brew install resvg' or 'cargo install resvg'")
	}

	// Create temp files
	tmpDir := os.TempDir()
	svgPath := filepath.Join(tmpDir, "horoscope_wheel.svg")
	pngPath := filepath.Join(tmpDir, "horoscope_wheel.png")

	// Write SVG to temp file
	if err := os.WriteFile(svgPath, svgData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write SVG: %w", err)
	}
	defer func() { _ = os.Remove(svgPath) }()
	defer func() { _ = os.Remove(pngPath) }()

	// Run resvg
	args := []string{
		"--width", fmt.Sprintf("%d", width),
		"--height", fmt.Sprintf("%d", height),
		svgPath,
		pngPath,
	}

	cmd := exec.Command(resvgPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("resvg failed: %w\noutput: %s", err, output)
	}

	// Check if PNG was created
	if _, err := os.Stat(pngPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("resvg did not create output file, output: %s", output)
	}

	// Read PNG data
	pngData, err := os.ReadFile(pngPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PNG: %w", err)
	}

	return pngData, nil
}

// HasResvg checks if resvg is installed
func HasResvg() bool {
	_, err := exec.LookPath("resvg")
	return err == nil
}

// ResvgVersion returns the installed resvg version
func ResvgVersion() string {
	cmd := exec.Command("resvg", "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}
