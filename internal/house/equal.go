package house

import (
	"time"

	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

// EqualCalculator implements the Equal House system
// Each house spans exactly 30 degrees starting from the Ascendant
type EqualCalculator struct{}

// Calculate computes house cusps using the Equal House system
func (e *EqualCalculator) Calculate(latitude, longitude float64, t time.Time) *Cusps {
	asc := position.CalculateAscendant(latitude, longitude, t)
	mc := position.CalculateMC(longitude, t)

	cusps := &Cusps{
		System:     Equal,
		Ascendant:  asc,
		MC:         mc,
		IC:         position.NormalizeAngle(mc + 180),
		Descendant: position.NormalizeAngle(asc + 180),
	}

	// Each house starts 30 degrees from the previous
	for i := 0; i < 12; i++ {
		cuspLon := position.NormalizeAngle(asc + float64(i)*30)
		cusps.Houses[i] = House{
			Number: i + 1,
			Cusp:   cuspLon,
			Sign:   horoscope.LongitudeToZodiac(cuspLon).Sign,
		}
	}

	return cusps
}
