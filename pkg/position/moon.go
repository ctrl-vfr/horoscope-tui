package position

import "math"

// calculateMoon computes the Moon's geocentric position with perturbations
func calculateMoon(d float64) Position {
	elem := PlanetElements[Moon].AtDay(d)

	// Also need Sun's mean anomaly for perturbations
	sunElem := PlanetElements[Sun].AtDay(d)
	Ms := DegreesToRadians(sunElem.M)
	Ls := DegreesToRadians(sunElem.M + sunElem.W)

	// Moon's mean longitude and mean anomaly
	N := DegreesToRadians(elem.N)
	I := DegreesToRadians(elem.I)
	Mm := DegreesToRadians(elem.M)
	Lm := N + DegreesToRadians(elem.W) + Mm

	// Moon's mean elongation
	D := Lm - Ls
	// Moon's argument of latitude
	F := Lm - N

	// Perturbations in longitude (degrees)
	pertLon := -1.274*math.Sin(Mm-2*D) + // Evection
		0.658*math.Sin(2*D) + // Variation
		-0.186*math.Sin(Ms) + // Yearly equation
		-0.059*math.Sin(2*Mm-2*D) +
		-0.057*math.Sin(Mm-2*D+Ms) +
		0.053*math.Sin(Mm+2*D) +
		0.046*math.Sin(2*D-Ms) +
		0.041*math.Sin(Mm-Ms) +
		-0.035*math.Sin(D) + // Parallactic equation
		-0.031*math.Sin(Mm+Ms) +
		-0.015*math.Sin(2*F-2*D) +
		0.011*math.Sin(Mm-4*D)

	// Perturbations in latitude (degrees)
	pertLat := -0.173*math.Sin(F-2*D) +
		-0.055*math.Sin(Mm-F-2*D) +
		-0.046*math.Sin(Mm+F-2*D) +
		0.033*math.Sin(F+2*D) +
		0.017*math.Sin(2*Mm+F)

	// Perturbations in distance (Earth radii)
	pertDist := -0.58*math.Cos(Mm-2*D) +
		-0.46*math.Cos(2*D)

	// Solve Kepler's equation
	E := solveKepler(Mm, elem.E)

	// Position in orbital plane
	xv := elem.A * (math.Cos(E) - elem.E)
	yv := elem.A * math.Sqrt(1-elem.E*elem.E) * math.Sin(E)

	// True anomaly and distance
	v := math.Atan2(yv, xv)
	r := math.Sqrt(xv*xv + yv*yv)

	// Apply distance perturbation
	r += pertDist

	// Ecliptic coordinates
	W := DegreesToRadians(elem.W)
	xh := r * (math.Cos(N)*math.Cos(v+W) - math.Sin(N)*math.Sin(v+W)*math.Cos(I))
	yh := r * (math.Sin(N)*math.Cos(v+W) + math.Cos(N)*math.Sin(v+W)*math.Cos(I))
	zh := r * math.Sin(v+W) * math.Sin(I)

	// Convert to ecliptic longitude/latitude
	lon := RadiansToDegrees(math.Atan2(yh, xh)) + pertLon
	lat := RadiansToDegrees(math.Atan2(zh, math.Sqrt(xh*xh+yh*yh))) + pertLat

	return Position{
		Body:              Moon,
		EclipticLongitude: NormalizeAngle(lon),
		EclipticLatitude:  lat,
		Distance:          r,
	}
}

// MoonPhase calculates the moon phase (0 = new, 0.5 = full)
func MoonPhase(d float64) float64 {
	sunPos := calculateSun(d)
	moonPos := calculateMoon(d)

	// Phase angle is difference in ecliptic longitude
	phase := moonPos.EclipticLongitude - sunPos.EclipticLongitude
	phase = NormalizeAngle(phase)

	return phase / 360.0
}

// MoonPhaseName returns the name of the current moon phase
func MoonPhaseName(d float64) string {
	phase := MoonPhase(d)
	switch {
	case phase < 0.0625 || phase >= 0.9375:
		return "New Moon"
	case phase < 0.1875:
		return "Waxing Crescent"
	case phase < 0.3125:
		return "First Quarter"
	case phase < 0.4375:
		return "Waxing Gibbous"
	case phase < 0.5625:
		return "Full Moon"
	case phase < 0.6875:
		return "Waning Gibbous"
	case phase < 0.8125:
		return "Last Quarter"
	default:
		return "Waning Crescent"
	}
}
