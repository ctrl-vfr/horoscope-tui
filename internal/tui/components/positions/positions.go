// Package positions provides the planetary positions table component.
package positions

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ctrl-vfr/astral-tui/internal/i18n"
	"github.com/ctrl-vfr/astral-tui/internal/tui/styles"
	"github.com/ctrl-vfr/astral-tui/pkg/horoscope"
	"github.com/ctrl-vfr/astral-tui/pkg/position"
)

// Model is the positions table component state.
type Model struct {
	viewport viewport.Model
	table    table.Model
	transits []position.Position
	natal    []position.Position
	chart    *horoscope.Chart
	width    int
	height   int
	focused  bool
}

// New creates a new positions model.
func New() Model {
	return Model{}
}

// Init initializes the positions component.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages for the positions component.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.focused {
		m.viewport, cmd = m.viewport.Update(msg)
	}
	return m, cmd
}

// SetSize sets the component dimensions.
func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	m.viewport = viewport.New(width-4, height-5)
	m.table = m.buildTable()
	m.viewport.SetContent(m.table.View())
	return m
}

// SetPositions sets the transit positions to display.
func (m Model) SetPositions(positions []position.Position) Model {
	m.transits = positions
	if m.width > 0 {
		m.table = m.buildTable()
		m.viewport.SetContent(m.table.View())
	}
	return m
}

// SetTransits updates the current transit positions.
func (m Model) SetTransits(positions []position.Position) Model {
	m.transits = positions
	if m.width > 0 {
		m.table = m.buildTable()
		m.viewport.SetContent(m.table.View())
	}
	return m
}

// SetChart sets the natal chart to display alongside transits.
func (m Model) SetChart(chart *horoscope.Chart) Model {
	m.chart = chart
	m.natal = chart.Positions
	if m.width > 0 {
		m.table = m.buildTable()
		m.viewport.SetContent(m.table.View())
	}
	return m
}

// SetFocus sets the focus state of the component.
func (m Model) SetFocus(focused bool) Model {
	m.focused = focused
	return m
}

func (m Model) buildTable() table.Model {
	var columns []table.Column
	var rows []table.Row

	// Available width for table content (minus borders and padding)
	availableWidth := m.width - 9
	availableWidth = max(availableWidth, 30)

	if m.chart != nil {
		fixedWidth := 12
		flexWidth := availableWidth - fixedWidth
		planetWidth := flexWidth / 4
		posWidth := (flexWidth - planetWidth) / 2

		columns = []table.Column{
			{Title: "", Width: 3},
			{Title: i18n.T("PositionPlanet"), Width: planetWidth},
			{Title: i18n.T("PositionNatal"), Width: posWidth},
			{Title: "℞", Width: 3},
			{Title: i18n.T("PositionTransit"), Width: posWidth},
			{Title: "℞", Width: 3},
		}
		rows = m.buildCombinedRows()
	} else {
		fixedWidth := 9
		flexWidth := availableWidth - fixedWidth
		planetWidth := flexWidth / 3
		posWidth := flexWidth - planetWidth

		columns = []table.Column{
			{Title: "", Width: 3},
			{Title: i18n.T("PositionPlanet"), Width: planetWidth},
			{Title: i18n.T("PositionPosition"), Width: posWidth},
			{Title: "℞", Width: 3},
		}
		rows = m.buildTransitRows()
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("94")).
		BorderBottom(true).
		Bold(true).
		Foreground(styles.ColorBright)
	s.Cell = s.Cell.Foreground(styles.ColorTextWarm)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(len(rows)+1),
		table.WithStyles(s),
	)
	t.Blur()

	return t
}

func (m Model) buildTransitRows() []table.Row {
	rows := make([]table.Row, 0, len(m.transits))
	for _, pos := range m.transits {
		zodiac := horoscope.LongitudeToZodiac(pos.EclipticLongitude)
		retro := ""
		if pos.Retrograde {
			retro = position.RetrogradeSymbol
		}
		rows = append(rows, table.Row{
			pos.Body.Symbol(),
			pos.Body.String(),
			fmt.Sprintf("%s %02d°%02d'", zodiac.Sign.Symbol(), zodiac.Degrees, zodiac.Minutes),
			retro,
		})
	}
	return rows
}

func (m Model) buildCombinedRows() []table.Row {
	rows := make([]table.Row, 0)

	transitMap := make(map[position.CelestialBody]position.Position)
	for _, pos := range m.transits {
		transitMap[pos.Body] = pos
	}

	for _, natal := range m.natal {
		natalZodiac := horoscope.LongitudeToZodiac(natal.EclipticLongitude)
		natalRetro := ""
		if natal.Retrograde {
			natalRetro = position.RetrogradeSymbol
		}

		transitPos := ""
		transitRetro := ""
		if transit, ok := transitMap[natal.Body]; ok {
			transitZodiac := horoscope.LongitudeToZodiac(transit.EclipticLongitude)
			transitPos = fmt.Sprintf("%s %02d°%02d'", transitZodiac.Sign.Symbol(), transitZodiac.Degrees, transitZodiac.Minutes)
			if transit.Retrograde {
				transitRetro = position.RetrogradeSymbol
			}
		}

		rows = append(rows, table.Row{
			natal.Body.Symbol(),
			natal.Body.String(),
			fmt.Sprintf("%s %02d°%02d'", natalZodiac.Sign.Symbol(), natalZodiac.Degrees, natalZodiac.Minutes),
			natalRetro,
			transitPos,
			transitRetro,
		})
	}
	return rows
}

// View renders the positions component.
func (m Model) View() string {
	borderColor := lipgloss.Color("94")
	if m.focused {
		borderColor = styles.ColorPrimary
	}

	title := i18n.T("PositionTransits")
	if m.chart != nil {
		title = i18n.T("PositionBoth")
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorBright).
		Render(title)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(m.width-2).
		Height(m.height-2).
		Padding(0, 1)

	return box.Render(header + "\n" + m.viewport.View())
}
