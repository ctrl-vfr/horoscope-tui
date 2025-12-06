package tui

import (
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/ctrl-vfr/astral-tui/internal/i18n"
	"github.com/ctrl-vfr/astral-tui/internal/tui/styles"
)

// View renders the main TUI view.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return i18n.T("StatusLoading")
	}
	return zone.Scan(m.buildView())
}

func (m Model) buildView() string {
	headerStyle := styles.PanelBorder.Width(m.width - 2)
	headerView := headerStyle.Render(m.header.View())

	leftCol := lipgloss.JoinVertical(lipgloss.Left,
		m.wheel.View(),
		m.positions.View(),
	)

	var rightCol string
	if m.chart != nil {
		rightCol = m.interp.View()
	} else {
		rightCol = m.form.View()
	}

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	statusLine := ""
	if m.status != "" {
		statusLine = "\n" + styles.DimStyle.Render(m.status)
	}

	helpLine := m.renderHelpLine()

	return lipgloss.JoinVertical(lipgloss.Left, headerView, mainContent) + statusLine + helpLine
}

func (m Model) renderHelpLine() string {
	keyStyle := lipgloss.NewStyle().Foreground(styles.ColorBright)
	sepStyle := styles.DimStyle

	help := keyStyle.Render("tab") + sepStyle.Render("/") + keyStyle.Render("shift+tab") + sepStyle.Render(i18n.T("NavNavigate")+" • ") +
		keyStyle.Render("↑↓") + sepStyle.Render(i18n.T("NavScroll")+" • ")

	if m.chart != nil {
		help += keyStyle.Render("esc") + sepStyle.Render(i18n.T("NavNewQuestion")+" • ")
	}

	help += keyStyle.Render("ctrl+c") + sepStyle.Render(i18n.T("NavQuit"))

	return "\n" + lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(help)
}

func (m Model) updateLayout() Model {
	headerHeight := 3
	helpLineHeight := 2
	contentHeight := m.height - headerHeight - helpLineHeight

	leftWidth := m.width / 2
	rightWidth := m.width - leftWidth

	wheelHeight := contentHeight * 60 / 100
	posHeight := contentHeight - wheelHeight

	m.header = m.header.SetSize(m.width)
	m.wheel = m.wheel.SetSize(leftWidth, wheelHeight)
	m.positions = m.positions.SetSize(leftWidth, posHeight)
	m.form = m.form.SetSize(rightWidth, contentHeight)
	m.interp = m.interp.SetSize(rightWidth, contentHeight)

	return m
}

func (m *Model) cycleFocus() {
	switch m.focus {
	case FocusForm:
		if m.chart != nil {
			m.focus = FocusInterp
		}
	case FocusInterp:
		if m.chart != nil {
			m.focus = FocusPositions
		}
	case FocusPositions:
		if !m.form.IsSubmitted() {
			m.focus = FocusForm
		} else {
			m.focus = FocusInterp
		}
	}
}

func (m *Model) cycleFocusReverse() {
	switch m.focus {
	case FocusForm:
		if m.chart != nil {
			m.focus = FocusPositions
		}
	case FocusInterp:
		if !m.form.IsSubmitted() {
			m.focus = FocusForm
		} else if m.chart != nil {
			m.focus = FocusPositions
		}
	case FocusPositions:
		m.focus = FocusInterp
	}
}

func (m Model) updateFocus() Model {
	m.interp = m.interp.SetFocus(m.focus == FocusInterp)
	m.positions = m.positions.SetFocus(m.focus == FocusPositions)
	return m
}
