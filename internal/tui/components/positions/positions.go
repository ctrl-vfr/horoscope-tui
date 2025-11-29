package positions

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

type Model struct {
	viewport  viewport.Model
	table     table.Model
	transits  []position.Position
	natal     []position.Position
	chart     *horoscope.Chart
	width     int
	height    int
	focused   bool
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.focused {
		m.viewport, cmd = m.viewport.Update(msg)
	}
	return m, cmd
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	m.viewport = viewport.New(width-4, height-5)
	m.table = m.buildTable()
	m.viewport.SetContent(m.table.View())
	return m
}

func (m Model) SetPositions(positions []position.Position) Model {
	// Always show current transits before chart validation
	m.transits = position.CalculateAll(time.Now())
	if m.width > 0 {
		m.table = m.buildTable()
		m.viewport.SetContent(m.table.View())
	}
	return m
}

func (m Model) SetChart(chart *horoscope.Chart) Model {
	m.chart = chart
	m.natal = chart.Positions
	m.transits = position.CalculateAll(time.Now())
	if m.width > 0 {
		m.table = m.buildTable()
		m.viewport.SetContent(m.table.View())
	}
	return m
}

func (m Model) SetFocus(focused bool) Model {
	m.focused = focused
	return m
}

func (m Model) buildTable() table.Model {
	var columns []table.Column
	var rows []table.Row

	if m.chart != nil {
		// Both natal and transits
		columns = []table.Column{
			{Title: "", Width: 2},
			{Title: "Planète", Width: 9},
			{Title: "Natal", Width: 11},
			{Title: "℞", Width: 2},
			{Title: "Transit", Width: 11},
			{Title: "℞", Width: 2},
		}
		rows = m.buildCombinedRows()
	} else {
		// Transits only
		columns = []table.Column{
			{Title: "", Width: 2},
			{Title: "Planète", Width: 10},
			{Title: "Position", Width: 12},
			{Title: "℞", Width: 2},
		}
		rows = m.buildTransitRows()
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("202"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	s.Cell = s.Cell.Foreground(lipgloss.Color("252"))

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

func (m Model) View() string {
	borderColor := lipgloss.Color("94")
	if m.focused {
		borderColor = lipgloss.Color("208")
	}

	title := "Transits"
	if m.chart != nil {
		title = "Natal / Transits"
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("202")).
		Render(title)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(m.width - 2).
		Height(m.height - 2).
		Padding(0, 1)

	return box.Render(header + "\n" + m.viewport.View())
}
