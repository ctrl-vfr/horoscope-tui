package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"

	"github.com/ctrl-vfr/horoscope-tui/internal/tui/components/form"
	"github.com/ctrl-vfr/horoscope-tui/internal/tui/components/header"
	"github.com/ctrl-vfr/horoscope-tui/internal/tui/components/interp"
	"github.com/ctrl-vfr/horoscope-tui/internal/tui/components/positions"
	"github.com/ctrl-vfr/horoscope-tui/internal/tui/components/wheel"
	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

// FocusArea represents the currently focused UI component.
type FocusArea int

const (
	FocusForm FocusArea = iota
	FocusInterp
	FocusPositions
)

// Model is the main TUI application model.
type Model struct {
	width  int
	height int

	header    header.Model
	form      form.Model
	wheel     wheel.Model
	interp    interp.Model
	positions positions.Model

	chart   *horoscope.Chart
	focus   FocusArea
	loading bool
	status  string
}

// NewModel creates a new TUI model with default state.
func NewModel() Model {
	zone.NewGlobal()

	today := time.Now()
	todayPositions := position.CalculateAll(today)

	return Model{
		header:    header.New().SetPositions(todayPositions),
		form:      form.New(),
		wheel:     wheel.New(),
		interp:    interp.New(),
		positions: positions.New().SetPositions(todayPositions),
		focus:     FocusForm,
	}
}

// Init initializes the TUI application.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.form.Init(),
		tea.EnterAltScreen,
		m.wheel.GenerateWheel(),
	)
}
