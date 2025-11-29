package interp

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/ctrl-vfr/horoscope-tui/internal/client"
	"github.com/ctrl-vfr/horoscope-tui/internal/tui/messages"
	"github.com/ctrl-vfr/horoscope-tui/internal/tui/styles"
	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
)

type Model struct {
	viewport viewport.Model
	content  string
	rendered string
	width    int
	height   int
	loading  bool
	complete bool
	focused  bool
	err      error
}

func New() Model {
	return Model{}
}

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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.InterpReadyMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.complete = true
			m.content = msg.Content
			m.renderMarkdown()
			m.viewport.SetContent(m.rendered)
		}

	case tea.KeyMsg:
		if m.focused {
			switch msg.String() {
			case "r":
				if m.err != nil {
					// Allow retry - handled by parent
				}
			}
		}
	}

	var cmd tea.Cmd
	if m.focused {
		m.viewport, cmd = m.viewport.Update(msg)
	}
	return m, cmd
}

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

func (m Model) SetFocus(focused bool) Model {
	m.focused = focused
	return m
}

func (m Model) StartInterpretation(chart *horoscope.Chart, userContext string) (Model, tea.Cmd) {
	m.loading = true
	m.complete = false
	m.content = ""
	m.err = nil
	return m, m.fetchCmd(chart, userContext)
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

	var header string
	if m.loading {
		header = headerStyle.Render("Interprétation") + " " + dimStyle.Render("(chargement...)")
	} else if m.err != nil {
		header = headerStyle.Render("Interprétation") + " " + lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Render("[Erreur]")
	} else {
		header = headerStyle.Render("Interprétation")
	}

	var content string
	if m.err != nil {
		content = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Render(m.err.Error())
	} else if m.content == "" && !m.loading {
		content = dimStyle.Render("En attente du thème natal...")
	} else {
		content = m.viewport.View()
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(m.width-2).
		Height(m.height-2).
		Padding(0, 1)

	return box.Render(header + "\n" + strings.Repeat("─", m.width-6) + "\n" + content)
}

func (m Model) IsLoading() bool {
	return m.loading
}

func (m Model) HasError() bool {
	return m.err != nil
}

func (m Model) Reset() Model {
	m.content = ""
	m.loading = false
	m.complete = false
	m.err = nil
	m.viewport.SetContent("")
	return m
}
