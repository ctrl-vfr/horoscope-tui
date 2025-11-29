package position

import "math"

// calculateNorthNode computes the Mean North Node (Rahu)
func calculateNorthNode(d float64) Position {
	// Mean North Node moves retrograde at ~19.35 degrees/year
	// Using Moon's ascending node
	elem := PlanetElements[Moon]
	N := elem.N + elem.NRate*d
	N = NormalizeAngle(N)

	return Position{
		Body:              NorthNode,
		EclipticLongitude: N,
		EclipticLatitude:  0,
		Distance:          0, // Not applicable
	}
}

// calculateSouthNode computes the South Node (Ketu)
// Always opposite to North Node
func calculateSouthNode(d float64) Position {
	northNode := calculateNorthNode(d)

	return Position{
		Body:              SouthNode,
		EclipticLongitude: NormalizeAngle(northNode.EclipticLongitude + 180),
		EclipticLatitude:  0,
		Distance:          0,
	}
}

// TrueNorthNode calculates the True North Node with nutation correction
func TrueNorthNode(d float64) float64 {
	meanNode := calculateNorthNode(d).EclipticLongitude

	// Simplified nutation correction
	// Full calculation would use IAU 2000A model
	sunElem := PlanetElements[Sun].AtDay(d)

	// Longitude of Moon's ascending node
	omega := DegreesToRadians(125.04 - 1934.136*d/36525)

	// Simplified nutation in longitude (arcseconds to degrees)
	deltaPsi := -17.20*math.Sin(omega)/3600 - 1.32*math.Sin(2*DegreesToRadians(sunElem.M+sunElem.W))/3600

	return NormalizeAngle(meanNode + deltaPsi)
}
