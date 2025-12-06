package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ctrl-vfr/horoscope-tui/internal/house"
	"github.com/ctrl-vfr/horoscope-tui/internal/i18n"
	"github.com/ctrl-vfr/horoscope-tui/internal/tui/components/wheel"
	"github.com/ctrl-vfr/horoscope-tui/internal/tui/messages"
	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

// Update handles messages for the main TUI model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.cycleFocus()
			m = m.updateFocus()
		case "shift+tab":
			m.cycleFocusReverse()
			m = m.updateFocus()
		case "esc":
			if m.chart != nil {
				m.form = m.form.Reset()
				m.interp = m.interp.Reset()
				m.chart = nil
				m.focus = FocusForm
				m = m.updateFocus()
				return m, m.form.Init()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.updateLayout()
		if m.wheel.IsReady() {
			cmds = append(cmds, m.wheel.RetransmitImage())
		}

	case messages.DateChangedMsg:
		positions := position.CalculateAll(msg.Date)
		m.header = m.header.SetPositions(positions)
		m.wheel = m.wheel.SetPositions(positions)
		m.positions = m.positions.SetPositions(positions)
		cmds = append(cmds, m.wheel.GenerateWheel())

	case messages.TransitDateChangedMsg:
		transitPositions := position.CalculateAll(msg.Date)
		m.positions = m.positions.SetTransits(transitPositions)

	case messages.GeocodingResultMsg:
		if msg.Err != nil {
			m.status = i18n.T("StatusGeocodingError") + msg.Err.Error()
			m.form = m.form.Reset()
		} else {
			m.status = i18n.T("StatusCalculating")
			dateTime, _ := m.form.GetDateTime()
			cmds = append(cmds, m.calculateChart(dateTime, msg.Latitude, msg.Longitude, msg.DisplayName))
		}

	case messages.ChartReadyMsg:
		m.chart = msg.Chart
		m.loading = false
		m.status = ""

		m.header = m.header.SetChart(m.chart.DateTime, m.chart.Location, m.chart.Positions)
		m.wheel = m.wheel.SetPositions(m.chart.Positions)
		m.positions = m.positions.SetChart(m.chart)

		// Set transit positions from form's transit date
		if transitDate, err := m.form.GetTransitDateTime(); err == nil {
			transitPositions := position.CalculateAll(transitDate)
			m.positions = m.positions.SetTransits(transitPositions)
		}

		cmds = append(cmds, m.wheel.GenerateWheel())

		userContext := m.form.GetUserContext()
		interpModel, interpCmd := m.interp.StartInterpretation(m.chart, userContext)
		m.interp = interpModel
		cmds = append(cmds, interpCmd)

		m.focus = FocusInterp
		m = m.updateFocus()

	case messages.WheelGeneratedMsg:
		var wheelCmd tea.Cmd
		m.wheel, wheelCmd = m.wheel.Update(msg)
		cmds = append(cmds, wheelCmd)

	case wheel.ImageTransmittedMsg:
		m.wheel, _ = m.wheel.Update(msg)

	case messages.InterpReadyMsg:
		var interpCmd tea.Cmd
		m.interp, interpCmd = m.interp.Update(msg)
		cmds = append(cmds, interpCmd)
	}

	// Update focused component
	switch m.focus {
	case FocusForm:
		if !m.form.IsSubmitted() {
			var formCmd tea.Cmd
			m.form, formCmd = m.form.Update(msg)
			cmds = append(cmds, formCmd)
		}
	case FocusInterp:
		var interpCmd tea.Cmd
		m.interp, interpCmd = m.interp.Update(msg)
		cmds = append(cmds, interpCmd)
	case FocusPositions:
		var posCmd tea.Cmd
		m.positions, posCmd = m.positions.Update(msg)
		cmds = append(cmds, posCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) calculateChart(dateTime time.Time, lat, lon float64, location string) tea.Cmd {
	return func() tea.Msg {
		positions := position.CalculateAll(dateTime)
		houseCusps := house.Calculate(lat, lon, dateTime)
		aspects := horoscope.CalculateAspects(positions, horoscope.DefaultOrbs)

		chart := &horoscope.Chart{
			DateTime:  dateTime,
			Latitude:  lat,
			Longitude: lon,
			Location:  location,
			Positions: positions,
			Houses:    houseCusps,
			Aspects:   aspects,
		}

		return messages.ChartReadyMsg{Chart: chart}
	}
}
