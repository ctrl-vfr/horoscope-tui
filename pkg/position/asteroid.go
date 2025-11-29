package position

// Asteroids and Chiron use the same calculation as planets
// Their orbital elements are defined in elements.go

// CalculateChiron computes Chiron's position
func CalculateChiron(d float64) Position {
	return calculatePlanet(Chiron, d)
}

// CalculateCeres computes Ceres' position
func CalculateCeres(d float64) Position {
	return calculatePlanet(Ceres, d)
}

// CalculatePallas computes Pallas' position
func CalculatePallas(d float64) Position {
	return calculatePlanet(Pallas, d)
}

// CalculateJuno computes Juno's position
func CalculateJuno(d float64) Position {
	return calculatePlanet(Juno, d)
}

// CalculateVesta computes Vesta's position
func CalculateVesta(d float64) Position {
	return calculatePlanet(Vesta, d)
}

// Asteroids returns all asteroid bodies
func Asteroids() []CelestialBody {
	return []CelestialBody{Chiron, Ceres, Pallas, Juno, Vesta}
}
