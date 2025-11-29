package house

import (
	"time"

	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

// System represents a house system type
type System int

const (
	Equal System = iota
	Placidus
)

func (s System) String() string {
	return []string{"Equal", "Placidus"}[s]
}

// House represents a single house in the chart
type House struct {
	Number int
	Cusp   float64 // Ecliptic longitude of house cusp
	Sign   horoscope.ZodiacSign
}

// Cusps contains all 12 house cusps plus angles
type Cusps struct {
	System     System
	Houses     [12]House
	Ascendant  float64
	MC         float64 // Midheaven (Medium Coeli)
	IC         float64 // Imum Coeli
	Descendant float64
}

// Calculator is the interface for house calculation systems
type Calculator interface {
	Calculate(latitude, longitude float64, t time.Time) *Cusps
}

// NewCalculator returns a house calculator for the given system
func NewCalculator(system System) Calculator {
	switch system {
	case Placidus:
		return &PlacidusCalculator{}
	default:
		return &EqualCalculator{}
	}
}

// GetHouse returns the house number (1-12) for a given ecliptic longitude
func (c *Cusps) GetHouse(longitude float64) int {
	longitude = position.NormalizeAngle(longitude)

	for i := 0; i < 12; i++ {
		nextHouse := (i + 1) % 12
		cusp := c.Houses[i].Cusp
		nextCusp := c.Houses[nextHouse].Cusp

		// Handle wrap-around at 0/360 degrees
		if nextCusp < cusp {
			if longitude >= cusp || longitude < nextCusp {
				return i + 1
			}
		} else {
			if longitude >= cusp && longitude < nextCusp {
				return i + 1
			}
		}
	}
	return 1
}
