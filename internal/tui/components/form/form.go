// Package form provides the birth data input form component.
package form

import (
	"errors"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/ctrl-vfr/horoscope-tui/internal/client"
	"github.com/ctrl-vfr/horoscope-tui/internal/i18n"
	"github.com/ctrl-vfr/horoscope-tui/internal/tui/messages"
	"github.com/ctrl-vfr/horoscope-tui/internal/tui/styles"
)

// Model is the form component state.
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

// New creates a new form model.
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
				Title(i18n.T("FormBirthDate")).
				Description(i18n.T("FormBirthDateDesc")).
				Placeholder(i18n.T("FormBirthDatePlaceholder")).
				Value(&m.dateStr).
				Validate(validateDate),
			huh.NewInput().
				Key("transit").
				Title(i18n.T("FormTransitDate")).
				Description(i18n.T("FormTransitDateDesc")).
				Placeholder(i18n.T("FormTransitDatePlaceholder")).
				Value(&m.transitDateStr).
				Validate(validateDate),
			huh.NewText().
				Key("context").
				Title(i18n.T("FormQuestion")).
				Description(i18n.T("FormQuestionDesc")).
				Placeholder(i18n.T("FormQuestionPlaceholder")).
				Value(&m.userContext).
				CharLimit(200),
		),
	).WithTheme(styles.HuhTheme()).
		WithShowHelp(false).
		WithShowErrors(true)
}

func validateDate(s string) error {
	if s == "" {
		return errors.New(i18n.T("ValidationRequired"))
	}
	_, err := time.Parse("02/01/2006", s)
	if err != nil {
		return errors.New(i18n.T("ValidationInvalidFormat"))
	}
	return nil
}

// Init initializes the form component.
func (m Model) Init() tea.Cmd {
	if m.missingCity {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the form component.
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

// SetSize sets the form dimensions.
func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

// IsSubmitted returns true if the form has been submitted.
func (m Model) IsSubmitted() bool {
	return m.submitted
}

// IsMissingCity returns true if the city environment variable is not set.
func (m Model) IsMissingCity() bool {
	return m.missingCity
}

// GetDateTime returns the parsed birth date and time.
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

// GetTransitDateTime returns the parsed transit date and time.
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

// GetUserContext returns the user's question or context.
func (m Model) GetUserContext() string {
	if m.form == nil {
		return ""
	}
	return m.form.GetString("context")
}

// View renders the form component.
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
			Render(i18n.T("StatusGeocoding")))
	}

	return box.Render(m.form.View())
}

func (m Model) missingCityContent() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	return errorStyle.Render(i18n.T("MissingCityError")) + "\n\n" +
		dimStyle.Render(i18n.T("MissingCityHint")) + "\n" +
		dimStyle.Render("export HOROSCOPE_CITY=\"Paris, France\"")
}

// Reset resets the form to its initial state.
func (m Model) Reset() Model {
	m.submitted = false
	m.loading = false
	m.err = nil
	if !m.missingCity {
		m.initForm()
	}
	return m
}
