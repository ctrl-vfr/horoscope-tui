package render

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/ansi/kitty"
)

// KittyGraphics handles image display in Kitty terminal
type KittyGraphics struct {
	writer io.Writer
}

// NewKittyGraphics creates a new Kitty graphics handler
func NewKittyGraphics() *KittyGraphics {
	return &KittyGraphics{
		writer: os.Stdout,
	}
}

// DisplayPNG displays a PNG image in the terminal (legacy method)
func (k *KittyGraphics) DisplayPNG(pngData []byte, cols, rows int) error {
	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		return fmt.Errorf("failed to decode PNG: %w", err)
	}

	opts := &kitty.Options{
		Action:       kitty.TransmitAndPut,
		Format:       kitty.PNG,
		Transmission: kitty.Direct,
		Columns:      cols,
		Rows:         rows,
		Chunk:        true,
		Quite:        0,
	}

	if err := kitty.EncodeGraphics(k.writer, img, opts); err != nil {
		return fmt.Errorf("failed to encode graphics: %w", err)
	}

	fmt.Fprintln(k.writer)
	return nil
}

// TransmitImage transmits image data to terminal with a unique ID for virtual placement
func TransmitImage(w io.Writer, pngData []byte, imageID int) error {
	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		return fmt.Errorf("failed to decode PNG: %w", err)
	}

	opts := &kitty.Options{
		Action:       kitty.Transmit,
		ID:           imageID,
		Format:       kitty.PNG,
		Transmission: kitty.Direct,
		Chunk:        true,
		Quite:        2, // Suppress all responses
	}

	return kitty.EncodeGraphics(w, img, opts)
}

// CreateVirtualPlacement creates a virtual placement for the image
func CreateVirtualPlacement(w io.Writer, imageID, cols, rows int) error {
	opts := &kitty.Options{
		Action:           kitty.Put,
		ID:               imageID,
		Columns:          cols,
		Rows:             rows,
		VirtualPlacement: true,
		Quite:            2,
	}

	return kitty.EncodeGraphics(w, nil, opts)
}

// BuildPlaceholderGrid builds a grid of Unicode placeholders for View()
// The terminal will replace these with the actual image pixels
func BuildPlaceholderGrid(imageID, cols, rows int) string {
	var sb strings.Builder

	for row := 0; row < rows; row++ {
		// Set foreground color to encode image ID (using 256-color palette)
		// Image ID must fit in the color value
		colorID := imageID % 256
		sb.WriteString(fmt.Sprintf("\x1b[38;5;%dm", colorID))

		for col := 0; col < cols; col++ {
			sb.WriteRune(kitty.Placeholder)
			sb.WriteRune(kitty.Diacritic(row))
			sb.WriteRune(kitty.Diacritic(col))
		}

		// Reset color and add newline
		sb.WriteString("\x1b[39m\n")
	}

	return sb.String()
}

// TransmitAndCreatePlacement transmits image and creates virtual placement in one call
func TransmitAndCreatePlacement(w io.Writer, pngData []byte, imageID, cols, rows int) error {
	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		return fmt.Errorf("failed to decode PNG: %w", err)
	}

	// First transmit the image
	transmitOpts := &kitty.Options{
		Action:       kitty.Transmit,
		ID:           imageID,
		Format:       kitty.PNG,
		Transmission: kitty.Direct,
		Chunk:        true,
		Quite:        2,
	}

	if err := kitty.EncodeGraphics(w, img, transmitOpts); err != nil {
		return fmt.Errorf("failed to transmit image: %w", err)
	}

	// Then create virtual placement
	placementOpts := &kitty.Options{
		Action:           kitty.Put,
		ID:               imageID,
		Columns:          cols,
		Rows:             rows,
		VirtualPlacement: true,
		Quite:            2,
	}

	if err := kitty.EncodeGraphics(w, nil, placementOpts); err != nil {
		return fmt.Errorf("failed to create placement: %w", err)
	}

	return nil
}

// DeleteImage deletes an image by ID from the terminal
func DeleteImage(w io.Writer, imageID int) {
	opts := &kitty.Options{
		Action: kitty.Delete,
		ID:     imageID,
		Delete: kitty.DeleteID,
		Quite:  2,
	}
	kitty.EncodeGraphics(w, nil, opts)
}

// PNGToImage converts PNG data to image.Image
func PNGToImage(pngData []byte) (image.Image, error) {
	return png.Decode(bytes.NewReader(pngData))
}

// SupportsKittyGraphics checks if the terminal supports Kitty graphics protocol
func SupportsKittyGraphics() bool {
	if os.Getenv("KITTY_WINDOW_ID") != "" {
		return true
	}

	term := os.Getenv("TERM")
	if term == "xterm-kitty" || term == "xterm-ghostty" {
		return true
	}

	termProgram := os.Getenv("TERM_PROGRAM")
	if termProgram == "kitty" || termProgram == "ghostty" || termProgram == "WezTerm" {
		return true
	}

	return false
}

// ClearImage clears any displayed images
func (k *KittyGraphics) ClearImage() {
	// a=d: delete
	// d=A: delete all images
	fmt.Fprint(k.writer, "\x1b_Ga=d,d=A\x1b\\")
}

// MoveCursor moves the cursor to a specific position
func (k *KittyGraphics) MoveCursor(row, col int) {
	fmt.Fprintf(k.writer, "\x1b[%d;%dH", row, col)
}

// SaveCursor saves the current cursor position
func (k *KittyGraphics) SaveCursor() {
	fmt.Fprint(k.writer, "\x1b7")
}

// RestoreCursor restores the saved cursor position
func (k *KittyGraphics) RestoreCursor() {
	fmt.Fprint(k.writer, "\x1b8")
}
