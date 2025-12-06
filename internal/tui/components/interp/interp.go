// Package interp provides the astrological interpretation display component.
package interp

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/ctrl-vfr/astral-tui/internal/client"
	"github.com/ctrl-vfr/astral-tui/internal/i18n"
	"github.com/ctrl-vfr/astral-tui/internal/tui/messages"
	"github.com/ctrl-vfr/astral-tui/internal/tui/styles"
	"github.com/ctrl-vfr/astral-tui/pkg/horoscope"
)

// Model is the interpretation component state.
type Model struct {
	viewport viewport.Model
	spinner  spinner.Model
	content  string
	rendered string
	question string
	width    int
	height   int
	loading  bool
	complete bool
	focused  bool
	err      error
}

// New creates a new interpretation model.
func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Moon
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	return Model{spinner: s}
}

// Init initializes the interpretation component.
func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) renderMarkdown() {
	if m.content == "" {
		m.rendered = ""
		return
	}
	wrapWidth := m.width - 6
	if wrapWidth < 20 {
		wrapWidth = 20
	}
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(styles.GlamourStyleJSON),
		glamour.WithWordWrap(wrapWidth),
	)
	if err != nil {
		m.rendered = m.content
		return
	}
	rendered, err := renderer.Render(m.content)
	if err != nil {
		m.rendered = m.content
		return
	}
	m.rendered = rendered
}

// Update handles messages for the interpretation component.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	if msg, ok := msg.(messages.InterpReadyMsg); ok {
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.complete = true
			m.content = msg.Content
			m.renderMarkdown()
			m.viewport.SetContent(m.rendered)
		}
	}

	// Update spinner while loading
	if m.loading {
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinnerCmd)
	}

	if m.focused {
		var vpCmd tea.Cmd
		m.viewport, vpCmd = m.viewport.Update(msg)
		cmds = append(cmds, vpCmd)
	}

	return m, tea.Batch(cmds...)
}

// SetSize sets the component dimensions.
func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	// viewport = total - border(2) - padding(2) - header(2)
	m.viewport = viewport.New(width-6, height-6)
	if m.content != "" {
		m.renderMarkdown()
		m.viewport.SetContent(m.rendered)
	}
	return m
}

// SetFocus sets the focus state of the component.
func (m Model) SetFocus(focused bool) Model {
	m.focused = focused
	return m
}

// StartInterpretation begins fetching an interpretation from OpenAI.
func (m Model) StartInterpretation(chart *horoscope.Chart, userContext string) (Model, tea.Cmd) {
	m.loading = true
	m.complete = false
	m.content = ""
	m.question = userContext
	m.err = nil
	return m, tea.Batch(m.fetchCmd(chart, userContext), m.spinner.Tick)
}

func (m Model) fetchCmd(chart *horoscope.Chart, userContext string) tea.Cmd {
	return func() tea.Msg {
		openaiClient, err := client.NewOpenAIClient()
		if err != nil {
			return messages.InterpReadyMsg{Err: err}
		}
		content, err := openaiClient.GetInterpretation(context.Background(), chart, userContext)
		return messages.InterpReadyMsg{Content: content, Err: err}
	}
}

// View renders the interpretation component.
func (m Model) View() string {
	borderColor := lipgloss.Color("94")
	if m.focused {
		borderColor = lipgloss.Color("208")
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("202"))

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	questionStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("223"))

	var header string
	switch {
	case m.loading:
		header = headerStyle.Render(i18n.T("InterpTitle")) + " " + m.spinner.View()
	case m.err != nil:
		header = headerStyle.Render(i18n.T("InterpTitle")) + " " + lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Render(i18n.T("InterpError"))
	default:
		header = headerStyle.Render(i18n.T("InterpTitle"))
	}

	var content string
	switch {
	case m.err != nil:
		content = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Render(m.err.Error())
	case m.loading:
		if m.question != "" {
			content = questionStyle.Render("« " + m.question + " »")
		}
	case m.content == "" && !m.loading:
		content = dimStyle.Render(i18n.T("StatusWaitingNatal"))
	default:
		if m.question != "" {
			content = questionStyle.Render("« "+m.question+" »") + "\n\n" + m.viewport.View()
		} else {
			content = m.viewport.View()
		}
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(m.width-2).
		Height(m.height-2).
		Padding(0, 1)

	return box.Render(header + "\n" + strings.Repeat("─", m.width-6) + "\n" + content)
}

// IsLoading returns true if an interpretation is being fetched.
func (m Model) IsLoading() bool {
	return m.loading
}

// HasError returns true if an error occurred.
func (m Model) HasError() bool {
	return m.err != nil
}

// Reset resets the component to its initial state.
func (m Model) Reset() Model {
	m.content = ""
	m.question = ""
	m.loading = false
	m.complete = false
	m.err = nil
	m.viewport.SetContent("")
	return m
}
