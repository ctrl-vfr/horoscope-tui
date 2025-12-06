// Package styles provides shared styling for TUI components.
package styles

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/ctrl-vfr/astral-tui/pkg/horoscope"
	"github.com/ctrl-vfr/astral-tui/pkg/position"
)

// Theme Colors
const (
	ColorPrimary   = lipgloss.Color("208") // Orange
	ColorSecondary = lipgloss.Color("160") // Red
	ColorAccent    = lipgloss.Color("53")  // Purple
	ColorText      = lipgloss.Color("255") // White
	ColorTextWarm  = lipgloss.Color("223") // Warm white
	ColorHighlight = lipgloss.Color("220") // Yellow
	ColorMuted     = lipgloss.Color("22")  // Dark green
	ColorBg        = lipgloss.Color("236") // Dark background
	ColorDim       = lipgloss.Color("240") // Dim gray
	ColorBright    = lipgloss.Color("202") // Bright orange
)

// Element colors
var (
	FireStyle  = lipgloss.NewStyle().Foreground(ColorPrimary)
	EarthStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("94"))
	AirStyle   = lipgloss.NewStyle().Foreground(ColorDim)
	WaterStyle = lipgloss.NewStyle().Foreground(ColorAccent)
)

// Planet colors
var (
	SunStyle     = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	MoonStyle    = lipgloss.NewStyle().Foreground(ColorTextWarm)
	MarsStyle    = lipgloss.NewStyle().Foreground(ColorSecondary)
	VenusStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	MercuryStyle = lipgloss.NewStyle().Foreground(ColorBright)
	JupiterStyle = lipgloss.NewStyle().Foreground(ColorAccent)
	SaturnStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("94"))
	UranusStyle  = lipgloss.NewStyle().Foreground(ColorMuted)
	NeptuneStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("54"))
	PlutoStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("52"))
)

// Aspect colors
var (
	HarmonicStyle = lipgloss.NewStyle().Foreground(ColorPrimary)
	TenseStyle    = lipgloss.NewStyle().Foreground(ColorSecondary)
	NeutralStyle  = lipgloss.NewStyle().Foreground(ColorHighlight)
)

// UI colors
var (
	HeaderStyle    = lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary)
	TitleStyle     = lipgloss.NewStyle().Bold(true).Foreground(ColorBright)
	BorderStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("94"))
	DimStyle       = lipgloss.NewStyle().Foreground(ColorDim)
	LabelStyle     = lipgloss.NewStyle().Foreground(ColorTextWarm)
	FocusedStyle   = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	UnfocusedStyle = lipgloss.NewStyle().Foreground(ColorDim)
)

// Border styles for panels
var (
	PanelBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("94"))

	FocusedBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary)
)

// StylePlanet returns a styled planet symbol
func StylePlanet(body position.CelestialBody) string {
	var style lipgloss.Style
	switch body {
	case position.Sun:
		style = SunStyle
	case position.Moon:
		style = MoonStyle
	case position.Mars:
		style = MarsStyle
	case position.Venus:
		style = VenusStyle
	case position.Mercury:
		style = MercuryStyle
	case position.Jupiter:
		style = JupiterStyle
	case position.Saturn:
		style = SaturnStyle
	case position.Uranus:
		style = UranusStyle
	case position.Neptune:
		style = NeptuneStyle
	case position.Pluto:
		style = PlutoStyle
	default:
		style = DimStyle
	}
	return style.Render(body.Symbol())
}

// StyleSign returns a styled zodiac sign
func StyleSign(sign horoscope.ZodiacSign) string {
	style := ElementStyle(sign.Element())
	return style.Render(sign.Symbol() + " " + sign.String()[:3])
}

// ElementStyle returns the style for an element
func ElementStyle(elem horoscope.Element) lipgloss.Style {
	switch elem {
	case horoscope.Fire:
		return FireStyle
	case horoscope.Earth:
		return EarthStyle
	case horoscope.Air:
		return AirStyle
	case horoscope.Water:
		return WaterStyle
	default:
		return lipgloss.NewStyle()
	}
}

// StyleAspect returns a styled aspect symbol
func StyleAspect(aspectType horoscope.AspectType) string {
	var style lipgloss.Style
	switch aspectType {
	case horoscope.Conjunction:
		style = NeutralStyle
	case horoscope.Sextile, horoscope.Trine:
		style = HarmonicStyle
	case horoscope.Square, horoscope.Opposition:
		style = TenseStyle
	default:
		style = DimStyle
	}
	return style.Render(aspectType.Symbol())
}

// HuhTheme returns an orange-themed huh form theme
func HuhTheme() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderForeground(ColorPrimary)
	t.Focused.Title = t.Focused.Title.Foreground(ColorBright)
	t.Focused.Description = t.Focused.Description.Foreground(ColorDim)
	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(ColorPrimary)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(ColorDim)

	t.Blurred.Base = t.Blurred.Base.BorderForeground(lipgloss.Color("94"))
	t.Blurred.Title = t.Blurred.Title.Foreground(ColorTextWarm)

	return t
}

// GlamourStyleJSON is the orange-themed glamour markdown style
var GlamourStyleJSON = []byte(`{
	"document": {
		"margin": 0
	},
	"heading": {
		"color": "202",
		"bold": true
	},
	"h1": {
		"prefix": "## ",
		"color": "208",
		"bold": true
	},
	"h2": {
		"prefix": "### ",
		"color": "202",
		"bold": true
	},
	"h3": {
		"prefix": "#### ",
		"color": "223"
	},
	"strong": {
		"color": "208",
		"bold": true
	},
	"emph": {
		"color": "223",
		"italic": true
	},
	"list": {
		"color": "223"
	},
	"item": {
		"color": "223"
	},
	"paragraph": {
		"color": "255"
	},
	"link": {
		"color": "208",
		"underline": true
	},
	"code": {
		"color": "202",
		"background_color": "236"
	},
	"code_block": {
		"color": "202",
		"margin": 1
	}
}`)
