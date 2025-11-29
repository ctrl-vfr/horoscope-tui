package render

import (
	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

// SVG color palette.
const (
	svgPrimary    = "#FF6B00"
	svgSecondary  = "#8B0000"
	svgAccent     = "#DC143C"
	svgPurple     = "#4B0082"
	svgPurpleDark = "#800080"
	svgTextLight  = "#F5E6D3"
	svgHighlight  = "#FFD700"
	svgMuted      = "#1A4D1A"
	svgBorder     = "#3D2D1F"
	svgDim        = "#4A4A4A"
	svgBright     = "#FF7518"
	svgTransit    = "#00CED1"
)

func getElementColor(elem horoscope.Element) string {
	switch elem {
	case horoscope.Fire:
		return svgPrimary
	case horoscope.Earth:
		return "#8B4513"
	case horoscope.Air:
		return svgDim
	case horoscope.Water:
		return svgPurpleDark
	default:
		return svgDim
	}
}

func getPlanetSVGColor(body position.CelestialBody) string {
	switch body {
	case position.Sun:
		return svgPrimary
	case position.Moon:
		return svgTextLight
	case position.Mercury:
		return svgBright
	case position.Venus:
		return "#FF69B4"
	case position.Mars:
		return svgAccent
	case position.Jupiter:
		return svgPurple
	case position.Saturn:
		return "#8B4513"
	case position.Uranus:
		return svgMuted
	case position.Neptune:
		return svgPurpleDark
	case position.Pluto:
		return svgSecondary
	default:
		return svgDim
	}
}

func getAspectColor(aspectType horoscope.AspectType) string {
	switch aspectType {
	case horoscope.Conjunction:
		return svgHighlight
	case horoscope.Sextile:
		return svgPrimary
	case horoscope.Trine:
		return svgBright
	case horoscope.Square:
		return svgAccent
	case horoscope.Opposition:
		return svgSecondary
	default:
		return svgDim
	}
}
