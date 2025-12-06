// Package position provides celestial body position calculations.
package position

// CelestialBody represents a celestial object.
type CelestialBody int

// RetrogradeSymbol is the astrological symbol for retrograde motion.
const RetrogradeSymbol = "℞"

// Celestial bodies used in astrological calculations.
const (
	Sun CelestialBody = iota
	Moon
	Mercury
	Venus
	Mars
	Jupiter
	Saturn
	Uranus
	Neptune
	Pluto
	NorthNode
	SouthNode
	Chiron
	Ceres
	Pallas
	Juno
	Vesta
)

// String returns the name of the celestial body
func (b CelestialBody) String() string {
	return bodyNames[b]
}

// Symbol returns the astrological symbol for the body
func (b CelestialBody) Symbol() string {
	return bodySymbols[b]
}

var bodyNames = map[CelestialBody]string{
	Sun:       "Sun",
	Moon:      "Moon",
	Mercury:   "Mercury",
	Venus:     "Venus",
	Mars:      "Mars",
	Jupiter:   "Jupiter",
	Saturn:    "Saturn",
	Uranus:    "Uranus",
	Neptune:   "Neptune",
	Pluto:     "Pluto",
	NorthNode: "North Node",
	SouthNode: "South Node",
	Chiron:    "Chiron",
	Ceres:     "Ceres",
	Pallas:    "Pallas",
	Juno:      "Juno",
	Vesta:     "Vesta",
}

var bodySymbols = map[CelestialBody]string{
	Sun:       "☉",
	Moon:      "☽",
	Mercury:   "☿",
	Venus:     "♀",
	Mars:      "♂",
	Jupiter:   "♃",
	Saturn:    "♄",
	Uranus:    "♅",
	Neptune:   "♆",
	Pluto:     "♇",
	NorthNode: "☊",
	SouthNode: "☋",
	Chiron:    "⚷",
	Ceres:     "⚳",
	Pallas:    "⚴",
	Juno:      "⚵",
	Vesta:     "⚶",
}

// AllBodies returns all celestial bodies in order
func AllBodies() []CelestialBody {
	return []CelestialBody{
		Sun, Moon, Mercury, Venus, Mars, Jupiter, Saturn,
		Uranus, Neptune, Pluto, NorthNode, SouthNode,
		Chiron, Ceres, Pallas, Juno, Vesta,
	}
}

// MainPlanets returns the traditional planets (Sun through Saturn)
func MainPlanets() []CelestialBody {
	return []CelestialBody{Sun, Moon, Mercury, Venus, Mars, Jupiter, Saturn}
}

// Position represents the calculated position of a celestial body
type Position struct {
	Body              CelestialBody
	EclipticLongitude float64 // degrees 0-360
	EclipticLatitude  float64 // degrees
	Distance          float64 // AU (or Earth radii for Moon)
	Retrograde        bool    // true if apparent retrograde motion
}

// CanBeRetrograde reports whether this body can exhibit retrograde motion.
func (b CelestialBody) CanBeRetrograde() bool {
	switch b {
	case Mercury, Venus, Mars, Jupiter, Saturn, Uranus, Neptune, Pluto, Chiron, Ceres, Pallas, Juno, Vesta:
		return true
	}
	return false
}

// IsMainPlanet reports whether this is a traditional planet (Sun through Saturn).
func (b CelestialBody) IsMainPlanet() bool {
	return b <= Saturn
}
