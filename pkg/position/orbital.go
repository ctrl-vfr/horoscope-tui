package position

import (
	"math"
	"time"
)

// Calculate computes the position of a celestial body at a given time
func Calculate(body CelestialBody, t time.Time) Position {
	d := DayNumber(t)
	return CalculateAtDay(body, d)
}

// CalculateAtDay computes the position for a day number from J2000
func CalculateAtDay(body CelestialBody, d float64) Position {
	switch body {
	case Moon:
		return calculateMoon(d)
	case NorthNode:
		return calculateNorthNode(d)
	case SouthNode:
		return calculateSouthNode(d)
	case Sun:
		return calculateSun(d)
	default:
		return calculatePlanet(body, d)
	}
}

// calculateSun computes the Sun's apparent position (geocentric)
func calculateSun(d float64) Position {
	elem := PlanetElements[Sun].AtDay(d)

	// Mean anomaly
	M := DegreesToRadians(elem.M)

	// Eccentric anomaly via Kepler's equation
	E := solveKepler(M, elem.E)

	// True anomaly and distance
	xv := elem.A * (math.Cos(E) - elem.E)
	yv := elem.A * math.Sqrt(1-elem.E*elem.E) * math.Sin(E)
	v := math.Atan2(yv, xv)
	r := math.Sqrt(xv*xv + yv*yv)

	// Sun's ecliptic longitude
	lonSun := RadiansToDegrees(v) + elem.W

	return Position{
		Body:              Sun,
		EclipticLongitude: NormalizeAngle(lonSun),
		EclipticLatitude:  0,
		Distance:          r,
	}
}

// calculatePlanet computes a planet's geocentric position
func calculatePlanet(body CelestialBody, d float64) Position {
	elem, ok := PlanetElements[body]
	if !ok {
		return Position{Body: body}
	}

	ce := elem.AtDay(d)

	// Solve Kepler's equation
	M := DegreesToRadians(ce.M)
	E := solveKepler(M, ce.E)

	// Position in orbital plane
	xv := ce.A * (math.Cos(E) - ce.E)
	yv := ce.A * math.Sqrt(1-ce.E*ce.E) * math.Sin(E)

	// True anomaly and distance
	v := math.Atan2(yv, xv)
	r := math.Sqrt(xv*xv + yv*yv)

	// Heliocentric coordinates in ecliptic plane
	N := DegreesToRadians(ce.N)
	I := DegreesToRadians(ce.I)
	W := DegreesToRadians(ce.W)

	xh := r * (math.Cos(N)*math.Cos(v+W) - math.Sin(N)*math.Sin(v+W)*math.Cos(I))
	yh := r * (math.Sin(N)*math.Cos(v+W) + math.Cos(N)*math.Sin(v+W)*math.Cos(I))
	zh := r * math.Sin(v+W) * math.Sin(I)

	// Get Sun's position for geocentric conversion
	sunPos := calculateSun(d)
	sunElem := PlanetElements[Sun].AtDay(d)
	sunM := DegreesToRadians(sunElem.M)
	sunE := solveKepler(sunM, sunElem.E)
	sunXv := sunElem.A * (math.Cos(sunE) - sunElem.E)
	sunYv := sunElem.A * math.Sqrt(1-sunElem.E*sunElem.E) * math.Sin(sunE)
	sunV := math.Atan2(sunYv, sunXv)
	sunR := math.Sqrt(sunXv*sunXv + sunYv*sunYv)
	sunLon := sunV + DegreesToRadians(sunElem.W)

	// Earth's heliocentric position (opposite of Sun's geocentric)
	xg := xh + sunR*math.Cos(sunLon)
	yg := yh + sunR*math.Sin(sunLon)
	zg := zh

	// Convert to ecliptic longitude/latitude
	lon := RadiansToDegrees(math.Atan2(yg, xg))
	lat := RadiansToDegrees(math.Atan2(zg, math.Sqrt(xg*xg+yg*yg)))
	dist := math.Sqrt(xg*xg + yg*yg + zg*zg)

	return Position{
		Body:              body,
		EclipticLongitude: NormalizeAngle(lon),
		EclipticLatitude:  lat,
		Distance:          dist + sunPos.Distance,
	}
}

// solveKepler solves Kepler's equation M = E - e*sin(E) iteratively
// meanAnom is mean anomaly in radians, e is eccentricity
// Returns eccentric anomaly E in radians
func solveKepler(meanAnom float64, e float64) float64 {
	// Initial approximation
	eccAnom := meanAnom + e*math.Sin(meanAnom)*(1.0+e*math.Cos(meanAnom))

	// Newton-Raphson iteration
	for i := 0; i < 15; i++ {
		next := eccAnom - (eccAnom-e*math.Sin(eccAnom)-meanAnom)/(1-e*math.Cos(eccAnom))
		if math.Abs(next-eccAnom) < 1e-12 {
			break
		}
		eccAnom = next
	}
	return eccAnom
}

// CalculateAll computes positions for all celestial bodies at a given time
func CalculateAll(t time.Time) []Position {
	d := DayNumber(t)
	bodies := AllBodies()
	positions := make([]Position, len(bodies))
	for i, body := range bodies {
		positions[i] = CalculateAtDay(body, d)
		if body.CanBeRetrograde() {
			positions[i].Retrograde = isRetrograde(body, d)
		}
	}
	return positions
}

// retrogradeCheckDelta is the number of days before/after to check motion direction.
const retrogradeCheckDelta = 1.0

// isRetrograde checks if a body is in retrograde motion at day d.
func isRetrograde(body CelestialBody, d float64) bool {
	posBefore := CalculateAtDay(body, d-retrogradeCheckDelta)
	posAfter := CalculateAtDay(body, d+retrogradeCheckDelta)
	motion := NormalizeMotion(posAfter.EclipticLongitude - posBefore.EclipticLongitude)
	return motion < 0
}

// CalculateAscendant computes the Ascendant (rising sign) for a location and time
func CalculateAscendant(latitude, longitude float64, t time.Time) float64 {
	jd := JulianDay(t)
	lst := LocalSiderealTime(jd, longitude)
	lstRad := DegreesToRadians(lst)
	latRad := DegreesToRadians(latitude)
	oblRad := DegreesToRadians(Obliquity)

	// Ascendant formula
	y := -math.Cos(lstRad)
	x := math.Sin(lstRad)*math.Cos(oblRad) + math.Tan(latRad)*math.Sin(oblRad)

	asc := RadiansToDegrees(math.Atan2(y, x))
	return NormalizeAngle(asc)
}

// CalculateMC computes the Midheaven (Medium Coeli)
func CalculateMC(longitude float64, t time.Time) float64 {
	jd := JulianDay(t)
	lst := LocalSiderealTime(jd, longitude)
	lstRad := DegreesToRadians(lst)
	oblRad := DegreesToRadians(Obliquity)

	mc := RadiansToDegrees(math.Atan2(math.Sin(lstRad), math.Cos(lstRad)*math.Cos(oblRad)))
	return NormalizeAngle(mc)
}
