// Package wheel provides the zodiac wheel visualization component.
package wheel

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ctrl-vfr/horoscope-tui/internal/i18n"
	"github.com/ctrl-vfr/horoscope-tui/internal/render"
	"github.com/ctrl-vfr/horoscope-tui/internal/tui/messages"
	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

const imageID = 42 // Fixed image ID for the zodiac wheel

// Model is the zodiac wheel component state.
type Model struct {
	positions        []position.Position
	pngData          []byte
	width            int
	height           int
	loading          bool
	imageReady       bool
	imageTransmitted bool
	cols             int
	rows             int
	err              error
}

// New creates a new wheel model.
func New() Model {
	return Model{}
}

// Init initializes the wheel component.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages for the wheel component.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.WheelGeneratedMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.pngData = msg.PNGData
			m.imageReady = true
			// Transmit image to terminal for virtual placement
			return m, m.transmitImage()
		}
	case ImageTransmittedMsg:
		m.imageTransmitted = true
	}
	return m, nil
}

// ImageTransmittedMsg signals that the image has been transmitted to terminal
type ImageTransmittedMsg struct{}

func (m Model) transmitImage() tea.Cmd {
	pngData := m.pngData
	cols := m.cols
	rows := m.rows
	return func() tea.Msg {
		// Delete any previous image with this ID
		render.DeleteImage(os.Stdout, imageID)
		// Transmit and create virtual placement
		if err := render.TransmitAndCreatePlacement(os.Stdout, pngData, imageID, cols, rows); err != nil {
			return messages.WheelGeneratedMsg{Err: err}
		}
		return ImageTransmittedMsg{}
	}
}

// SetSize sets the component dimensions.
func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	// Calculate cols and rows for the image placeholder
	// Leave margin for the border (2 for left/right, 2 for top/bottom)
	m.cols = width - 4
	m.rows = height - 4
	if m.cols < 10 {
		m.cols = 10
	}
	if m.rows < 5 {
		m.rows = 5
	}
	// Reset transmission state when size changes
	m.imageTransmitted = false
	return m
}

// SetPositions sets the natal positions for the wheel.
func (m Model) SetPositions(positions []position.Position) Model {
	m.positions = positions
	m.loading = true
	m.imageReady = false
	m.imageTransmitted = false
	return m
}

// HasPositions returns true if natal positions are set.
func (m Model) HasPositions() bool {
	return len(m.positions) > 0
}

// GenerateWheel generates the zodiac wheel image.
func (m Model) GenerateWheel() tea.Cmd {
	natalPositions := m.positions
	return func() tea.Msg {
		svgSize := 600
		generator := render.NewSVGWheelGenerator(svgSize)

		// Calculate today's transits
		transitPositions := position.CalculateAll(time.Now())

		var svgData []byte
		if len(natalPositions) == 0 {
			// No natal data yet - show transits only (on inner ring)
			svgData = generator.Generate(transitPositions)
		} else {
			// Show both natal (inner) and transits (outer)
			svgData = generator.GenerateWithTransits(natalPositions, transitPositions)
		}

		pngData, err := render.SVGToPNG(svgData, svgSize, svgSize)
		if err != nil {
			return messages.WheelGeneratedMsg{Err: err}
		}

		return messages.WheelGeneratedMsg{PNGData: pngData}
	}
}

// GenerateTransitsOnly generates a wheel with only today's transits
func (m Model) GenerateTransitsOnly() tea.Cmd {
	return func() tea.Msg {
		svgSize := 600
		generator := render.NewSVGWheelGenerator(svgSize)

		transitPositions := position.CalculateAll(time.Now())
		svgData := generator.Generate(transitPositions)

		pngData, err := render.SVGToPNG(svgData, svgSize, svgSize)
		if err != nil {
			return messages.WheelGeneratedMsg{Err: err}
		}

		return messages.WheelGeneratedMsg{PNGData: pngData}
	}
}

// View renders the wheel component.
func (m Model) View() string {
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("94")).
		Width(m.width-2).
		Height(m.height-2).
		Align(lipgloss.Center, lipgloss.Center)

	if m.loading {
		return borderStyle.Render(i18n.T("StatusGeneratingWheel"))
	}

	if m.err != nil {
		return borderStyle.Render(fmt.Sprintf("%s%v", i18n.T("StatusError"), m.err))
	}

	if !m.imageReady || len(m.pngData) == 0 {
		return borderStyle.Render(i18n.T("StatusWaitingData"))
	}

	// Check if terminal supports Kitty graphics
	if !render.SupportsKittyGraphics() || !render.HasResvg() {
		return m.asciiPlaceholder()
	}

	// Wait for image to be transmitted before showing placeholders
	if !m.imageTransmitted {
		return borderStyle.Render(i18n.T("StatusTransmittingImage"))
	}

	// Return Unicode placeholder grid wrapped in border
	grid := render.BuildPlaceholderGrid(imageID, m.cols, m.rows)
	return borderStyle.Render(grid)
}

func (m Model) asciiPlaceholder() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Border(lipgloss.RoundedBorder()).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("240"))

	return style.Render(i18n.T("WheelPlaceholder"))
}

// IsReady returns true if the wheel image is ready to display.
func (m Model) IsReady() bool {
	return m.imageReady && len(m.pngData) > 0
}

// RetransmitImage retransmits the wheel image to the terminal.
func (m Model) RetransmitImage() tea.Cmd {
	return m.transmitImage()
}
