package form

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/ctrl-vfr/horoscope-tui/internal/client"
	"github.com/ctrl-vfr/horoscope-tui/internal/tui/messages"
	"github.com/ctrl-vfr/horoscope-tui/internal/tui/styles"
)

type Model struct {
	form                 *huh.Form
	dateStr              string
	lastValidDate        string
	transitDateStr       string
	transitLastValidDate string
	userContext          string
	city                 string
	width                int
	height               int
	submitted            bool
	loading              bool
	err                  error
	missingCity          bool
}

func New() Model {
	city := os.Getenv("HOROSCOPE_CITY")
	today := time.Now().Format("02/01/2006")
	m := Model{
		dateStr:        today,
		transitDateStr: today,
		city:           city,
		missingCity:    city == "",
	}
	if !m.missingCity {
		m.initForm()
	}
	return m
}

func (m *Model) initForm() {
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("date").
				Title("Date de naissance").
				Description("Format: JJ/MM/AAAA").
				Placeholder("21/03/1990").
				Value(&m.dateStr).
				Validate(validateDate),
			huh.NewInput().
				Key("transit").
				Title("Date de transit").
				Description("Pour les prédictions (JJ/MM/AAAA)").
				Placeholder("01/01/2025").
				Value(&m.transitDateStr).
				Validate(validateDate),
			huh.NewText().
				Key("context").
				Title("Pose ta question à l'oracle").
				Description("Ex: Dois-je accepter ce job? C'est le bon moment pour...?").
				Placeholder("Qu'est-ce qui te tracasse?").
				Value(&m.userContext).
				CharLimit(200),
		),
	).WithTheme(styles.HuhTheme()).
		WithShowHelp(false).
		WithShowErrors(true)
}

func validateDate(s string) error {
	if s == "" {
		return fmt.Errorf("date requise")
	}
	_, err := time.Parse("02/01/2006", s)
	if err != nil {
		return fmt.Errorf("format invalide (JJ/MM/AAAA)")
	}
	return nil
}

func (m Model) Init() tea.Cmd {
	if m.missingCity {
		return nil
	}
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.missingCity || m.submitted {
		return m, nil
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	var cmds []tea.Cmd
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// Check if birth date changed and is valid
	if m.dateStr != m.lastValidDate && validateDate(m.dateStr) == nil {
		m.lastValidDate = m.dateStr
		if dateTime, err := m.GetDateTime(); err == nil {
			cmds = append(cmds, func() tea.Msg {
				return messages.DateChangedMsg{Date: dateTime}
			})
		}
	}

	// Check if transit date changed and is valid
	if m.transitDateStr != m.transitLastValidDate && validateDate(m.transitDateStr) == nil {
		m.transitLastValidDate = m.transitDateStr
		if transitTime, err := m.GetTransitDateTime(); err == nil {
			cmds = append(cmds, func() tea.Msg {
				return messages.TransitDateChangedMsg{Date: transitTime}
			})
		}
	}

	if m.form.State == huh.StateCompleted && !m.submitted {
		m.submitted = true
		m.loading = true
		cmds = append(cmds, m.geocodeCity())
	}

	return m, tea.Batch(cmds...)
}

func (m Model) geocodeCity() tea.Cmd {
	city := m.city
	return func() tea.Msg {
		geocoder := client.NewGeocodingClient()
		result, err := geocoder.Search(city)
		if err != nil {
			return messages.GeocodingResultMsg{Err: err}
		}
		return messages.GeocodingResultMsg{
			Latitude:    result.Latitude,
			Longitude:   result.Longitude,
			DisplayName: result.DisplayName,
		}
	}
}

func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

func (m Model) IsSubmitted() bool {
	return m.submitted
}

func (m Model) IsMissingCity() bool {
	return m.missingCity
}

func (m Model) GetDateTime() (time.Time, error) {
	dateStr := m.dateStr
	if m.form != nil {
		dateStr = m.form.GetString("date")
	}
	date, err := time.Parse("02/01/2006", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	now := time.Now()
	return time.Date(date.Year(), date.Month(), date.Day(),
		now.Hour(), now.Minute(), 0, 0, time.Local), nil
}

func (m Model) GetTransitDateTime() (time.Time, error) {
	dateStr := m.transitDateStr
	if m.form != nil {
		dateStr = m.form.GetString("transit")
	}
	date, err := time.Parse("02/01/2006", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(),
		12, 0, 0, 0, time.Local), nil
}

func (m Model) GetUserContext() string {
	if m.form == nil {
		return ""
	}
	return m.form.GetString("context")
}

func (m Model) View() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("94")).
		Width(m.width-2).
		Height(m.height-2).
		Padding(0, 1)

	if m.missingCity {
		return box.Render(m.missingCityContent())
	}

	if m.loading {
		return box.Render(lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("Recherche des coordonnées..."))
	}

	return box.Render(m.form.View())
}

func (m Model) missingCityContent() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	return errorStyle.Render("Variable HOROSCOPE_CITY manquante") + "\n\n" +
		dimStyle.Render("Définissez la variable d'environnement:") + "\n" +
		dimStyle.Render("export HOROSCOPE_CITY=\"Paris, France\"")
}

func (m Model) Reset() Model {
	m.submitted = false
	m.loading = false
	m.err = nil
	if !m.missingCity {
		m.initForm()
	}
	return m
}
