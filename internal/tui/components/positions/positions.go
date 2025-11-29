package positions

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ctrl-vfr/horoscope-tui/internal/tui/styles"
	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

type Model struct {
	viewport  viewport.Model
	positions []position.Position
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
	// viewport = total - border(2) - header(2)
	m.viewport = viewport.New(width-4, height-5)
	m.viewport.SetContent(m.buildContent())
	return m
}

func (m Model) SetPositions(positions []position.Position) Model {
	m.positions = positions
	if m.width > 0 {
		m.updateContent()
	}
	return m
}

func (m Model) SetChart(chart *horoscope.Chart) Model {
	m.chart = chart
	m.positions = chart.Positions
	m.updateContent()
	return m
}

func (m Model) SetFocus(focused bool) Model {
	m.focused = focused
	return m
}

func (m *Model) updateContent() {
	m.viewport.SetContent(m.buildContent())
}

func (m Model) buildContent() string {
	if len(m.positions) == 0 {
		return "En attente..."
	}

	var sb strings.Builder

	for _, pos := range m.positions {
		zodiac := horoscope.LongitudeToZodiac(pos.EclipticLongitude)
		retro := " "
		if pos.Retrograde {
			retro = position.RetrogradeSymbol
		}
		line := fmt.Sprintf(" %s %-8s %02d°%02d' %s %s\n",
			styles.StylePlanet(pos.Body),
			pos.Body.String(),
			zodiac.Degrees,
			zodiac.Minutes,
			styles.StyleSign(zodiac.Sign),
			retro)
		sb.WriteString(line)
	}

	// Element distribution
	sb.WriteString("\n")
	elements := m.countElements()
	sb.WriteString(fmt.Sprintf(" %s %s\n %s %s\n",
		styles.FireStyle.Render(fmt.Sprintf("Feu:%d", elements[horoscope.Fire])),
		styles.EarthStyle.Render(fmt.Sprintf("Terre:%d", elements[horoscope.Earth])),
		styles.AirStyle.Render(fmt.Sprintf("Air:%d", elements[horoscope.Air])),
		styles.WaterStyle.Render(fmt.Sprintf("Eau:%d", elements[horoscope.Water]))))

	return sb.String()
}

func (m Model) countElements() map[horoscope.Element]int {
	result := make(map[horoscope.Element]int)
	for _, pos := range m.positions {
		zodiac := horoscope.LongitudeToZodiac(pos.EclipticLongitude)
		result[zodiac.Sign.Element()]++
	}
	return result
}

func (m Model) View() string {
	borderColor := lipgloss.Color("94")
	if m.focused {
		borderColor = lipgloss.Color("208")
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("202")).
		Render("Positions")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(m.width-2).
		Height(m.height-2).
		Padding(0, 1)

	return box.Render(header + "\n" + m.viewport.View())
}
