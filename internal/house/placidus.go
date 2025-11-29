package house

import (
	"math"
	"time"

	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

// PlacidusCalculator implements the Placidus house system
// This is the most commonly used house system in Western astrology
type PlacidusCalculator struct{}

// Calculate computes house cusps using the Placidus system
func (p *PlacidusCalculator) Calculate(latitude, longitude float64, t time.Time) *Cusps {
	asc := position.CalculateAscendant(latitude, longitude, t)
	mc := position.CalculateMC(longitude, t)

	jd := position.JulianDay(t)
	lst := position.LocalSiderealTime(jd, longitude)
	ramc := lst // Right Ascension of MC in degrees

	latRad := position.DegreesToRadians(latitude)
	oblRad := position.DegreesToRadians(position.Obliquity)

	cusps := &Cusps{
		System:     Placidus,
		Ascendant:  asc,
		MC:         mc,
		IC:         position.NormalizeAngle(mc + 180),
		Descendant: position.NormalizeAngle(asc + 180),
	}

	// House 10 (MC) and House 4 (IC)
	cusps.Houses[9] = House{Number: 10, Cusp: mc, Sign: horoscope.LongitudeToZodiac(mc).Sign}
	cusps.Houses[3] = House{Number: 4, Cusp: position.NormalizeAngle(mc + 180),
		Sign: horoscope.LongitudeToZodiac(position.NormalizeAngle(mc + 180)).Sign}

	// House 1 (Ascendant) and House 7 (Descendant)
	cusps.Houses[0] = House{Number: 1, Cusp: asc, Sign: horoscope.LongitudeToZodiac(asc).Sign}
	cusps.Houses[6] = House{Number: 7, Cusp: position.NormalizeAngle(asc + 180),
		Sign: horoscope.LongitudeToZodiac(position.NormalizeAngle(asc + 180)).Sign}

	// Calculate intermediate cusps using Placidus method
	// Houses 11, 12, 2, 3 and their opposites 5, 6, 8, 9

	// House 11: 1/3 of semi-arc from MC
	cusps.Houses[10] = p.calculateCusp(11, ramc, 30, latRad, oblRad)
	// House 12: 2/3 of semi-arc from MC
	cusps.Houses[11] = p.calculateCusp(12, ramc, 60, latRad, oblRad)
	// House 2: 2/3 of semi-arc from IC
	cusps.Houses[1] = p.calculateCusp(2, ramc, 120, latRad, oblRad)
	// House 3: 1/3 of semi-arc from IC
	cusps.Houses[2] = p.calculateCusp(3, ramc, 150, latRad, oblRad)

	// Opposite houses
	cusps.Houses[4] = House{Number: 5, Cusp: position.NormalizeAngle(cusps.Houses[10].Cusp + 180),
		Sign: horoscope.LongitudeToZodiac(position.NormalizeAngle(cusps.Houses[10].Cusp + 180)).Sign}
	cusps.Houses[5] = House{Number: 6, Cusp: position.NormalizeAngle(cusps.Houses[11].Cusp + 180),
		Sign: horoscope.LongitudeToZodiac(position.NormalizeAngle(cusps.Houses[11].Cusp + 180)).Sign}
	cusps.Houses[7] = House{Number: 8, Cusp: position.NormalizeAngle(cusps.Houses[1].Cusp + 180),
		Sign: horoscope.LongitudeToZodiac(position.NormalizeAngle(cusps.Houses[1].Cusp + 180)).Sign}
	cusps.Houses[8] = House{Number: 9, Cusp: position.NormalizeAngle(cusps.Houses[2].Cusp + 180),
		Sign: horoscope.LongitudeToZodiac(position.NormalizeAngle(cusps.Houses[2].Cusp + 180)).Sign}

	return cusps
}

// calculateCusp calculates an intermediate house cusp using Placidus method
func (p *PlacidusCalculator) calculateCusp(houseNum int, ramc, offset, lat, obl float64) House {
	ra := position.NormalizeAngle(ramc + offset)
	raRad := position.DegreesToRadians(ra)

	// Iterative calculation for Placidus cusps
	var lon float64
	for i := 0; i < 20; i++ {
		// Calculate declination of the point
		decl := math.Asin(math.Sin(obl) * math.Sin(raRad))

		// Calculate semi-arc
		h := math.Acos(-math.Tan(lat) * math.Tan(decl))
		if math.IsNaN(h) {
			// Fall back to equal houses for extreme latitudes
			lon = ra
			break
		}

		// Calculate fraction based on house
		var f float64
		switch houseNum {
		case 11:
			f = 1.0 / 3.0
		case 12:
			f = 2.0 / 3.0
		case 2:
			f = 2.0 / 3.0
		case 3:
			f = 1.0 / 3.0
		}

		// Adjust RA
		newRA := ramc + offset
		if houseNum <= 3 || houseNum >= 10 {
			newRA = ramc + f*position.RadiansToDegrees(h)*2
		}

		raRad = position.DegreesToRadians(position.NormalizeAngle(newRA))

		// Convert RA to ecliptic longitude
		lon = position.RadiansToDegrees(math.Atan2(
			math.Sin(raRad)*math.Cos(obl)+math.Tan(decl)*math.Sin(obl),
			math.Cos(raRad),
		))
	}

	lon = position.NormalizeAngle(lon)
	return House{
		Number: houseNum,
		Cusp:   lon,
		Sign:   horoscope.LongitudeToZodiac(lon).Sign,
	}
}
