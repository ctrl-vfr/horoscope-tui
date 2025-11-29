package render

import (
	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

// SVG color palette.
const (
	svgPrimary   = "#FF6B00" // Orange
	svgAccent    = "#DC143C" // Crimson
	svgPurple    = "#4B0082" // Indigo
	svgTextLight = "#F5E6D3" // Beige
	svgBorder    = "#3D2D1F" // Dark brown
	svgBright    = "#FF7518" // Orange-red
	svgFire      = "#FF4500" // Orange-red
	svgEarth     = "#32CD32" // Lime green
	svgAir       = "#FFD700" // Gold
	svgWater     = "#1E90FF" // Dodger blue
)

func getElementColor(elem horoscope.Element) string {
	switch elem {
	case horoscope.Fire:
		return svgFire
	case horoscope.Earth:
		return svgEarth
	case horoscope.Air:
		return svgAir
	case horoscope.Water:
		return svgWater
	default:
		return svgTextLight
	}
}

func getPlanetSVGColor(body position.CelestialBody) string {
	switch body {
	case position.Sun:
		return "#FFD700" // Gold
	case position.Moon:
		return "#E6E6FA" // Lavender (silvery)
	case position.Mercury:
		return "#FFA500" // Orange
	case position.Venus:
		return "#FF69B4" // Hot pink
	case position.Mars:
		return "#FF4444" // Bright red
	case position.Jupiter:
		return "#9370DB" // Medium purple
	case position.Saturn:
		return "#DAA520" // Goldenrod
	case position.Uranus:
		return "#00CED1" // Dark turquoise
	case position.Neptune:
		return "#4169E1" // Royal blue
	case position.Pluto:
		return "#DC143C" // Crimson
	default:
		return svgTextLight
	}
}

func getAspectColor(aspectType horoscope.AspectType) string {
	switch aspectType {
	case horoscope.Conjunction:
		return "#FFD700" // Gold
	case horoscope.Sextile:
		return "#32CD32" // Lime green (harmonious)
	case horoscope.Trine:
		return "#00BFFF" // Deep sky blue (harmonious)
	case horoscope.Square:
		return "#FF4500" // Orange-red (tension)
	case horoscope.Opposition:
		return "#FF1493" // Deep pink (tension)
	default:
		return svgTextLight
	}
}
