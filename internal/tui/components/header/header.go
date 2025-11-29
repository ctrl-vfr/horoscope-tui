package header

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

type Model struct {
	width    int
	dateTime time.Time
	location string
	hasChart bool
	elements map[horoscope.Element]int
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) SetSize(width int) Model {
	m.width = width
	return m
}

func (m Model) SetPositions(positions []position.Position) Model {
	m.elements = calculateElements(positions)
	return m
}

func (m Model) SetChart(dateTime time.Time, location string, positions []position.Position) Model {
	m.dateTime = dateTime
	m.location = location
	m.hasChart = true
	m.elements = calculateElements(positions)
	return m
}

func calculateElements(positions []position.Position) map[horoscope.Element]int {
	elements := make(map[horoscope.Element]int)
	for _, pos := range positions {
		if pos.Body > position.Pluto {
			continue
		}
		zodiac := horoscope.LongitudeToZodiac(pos.EclipticLongitude)
		elements[zodiac.Sign.Element()]++
	}
	return elements
}

func (m Model) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("208"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("223"))

	fireStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	earthStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("94"))
	airStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	waterStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99"))

	// Title + location + elements on same line
	left := titleStyle.Render("✧ MON ORACLE ✧")
	if m.hasChart {
		dateStr := m.dateTime.Format("02 Jan 2006 15:04")
		left = left + "  " + dimStyle.Render(dateStr)
		if m.location != "" {
			left = left + "  " + dimStyle.Render("• "+m.location)
		}
	}

	right := ""
	if m.elements != nil {
		right = fmt.Sprintf("%s %s %s %s",
			fireStyle.Render(fmt.Sprintf("🔥%d", m.elements[horoscope.Fire])),
			earthStyle.Render(fmt.Sprintf("🪨%d", m.elements[horoscope.Earth])),
			airStyle.Render(fmt.Sprintf("💨%d", m.elements[horoscope.Air])),
			waterStyle.Render(fmt.Sprintf("💦%d", m.elements[horoscope.Water])))
	}

	width := m.width - 4
	if width < 40 {
		width = 76
	}

	return leftRightPad(left, right, width)
}

func leftRightPad(left, right string, width int) string {
	leftLen := lipgloss.Width(left)
	rightLen := lipgloss.Width(right)
	padding := width - leftLen - rightLen
	if padding < 1 {
		padding = 1
	}
	return left + strings.Repeat(" ", padding) + right
}
