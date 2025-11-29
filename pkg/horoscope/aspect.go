package horoscope

import (
	"math"

	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

// AspectType represents a type of astrological aspect
type AspectType int

const (
	Conjunction AspectType = iota // 0 degrees
	Sextile                       // 60 degrees
	Square                        // 90 degrees
	Trine                         // 120 degrees
	Opposition                    // 180 degrees
)

// String returns the name of the aspect
func (a AspectType) String() string {
	return aspectNames[a]
}

// Symbol returns the astrological symbol for the aspect
func (a AspectType) Symbol() string {
	return aspectSymbols[a]
}

// Angle returns the exact angle for this aspect type
func (a AspectType) Angle() float64 {
	return aspectAngles[a]
}

// IsHarmonic returns true if this is a harmonious aspect
func (a AspectType) IsHarmonic() bool {
	return a == Conjunction || a == Sextile || a == Trine
}

var aspectNames = map[AspectType]string{
	Conjunction: "Conjunction",
	Sextile:     "Sextile",
	Square:      "Square",
	Trine:       "Trine",
	Opposition:  "Opposition",
}

var aspectSymbols = map[AspectType]string{
	Conjunction: "☌",
	Sextile:     "⚹",
	Square:      "□",
	Trine:       "△",
	Opposition:  "☍",
}

var aspectAngles = map[AspectType]float64{
	Conjunction: 0,
	Sextile:     60,
	Square:      90,
	Trine:       120,
	Opposition:  180,
}

// Orbs contains the allowed orb (deviation) for each aspect type
type Orbs map[AspectType]float64

// DefaultOrbs provides standard orb values for aspects
var DefaultOrbs = Orbs{
	Conjunction: 8.0,
	Sextile:     6.0,
	Square:      7.0,
	Trine:       8.0,
	Opposition:  8.0,
}

// TightOrbs provides stricter orb values
var TightOrbs = Orbs{
	Conjunction: 5.0,
	Sextile:     4.0,
	Square:      5.0,
	Trine:       5.0,
	Opposition:  5.0,
}

// Aspect represents an aspect between two celestial bodies
type Aspect struct {
	Body1    position.CelestialBody
	Body2    position.CelestialBody
	Type     AspectType
	Angle    float64 // Actual angle between bodies
	Orb      float64 // Deviation from exact aspect
	Applying bool    // True if aspect is applying (getting tighter)
}

// CalculateAspects finds all aspects between a set of positions
func CalculateAspects(positions []position.Position, orbs Orbs) []Aspect {
	var aspects []Aspect

	for i := 0; i < len(positions); i++ {
		for j := i + 1; j < len(positions); j++ {
			if aspect := findAspect(positions[i], positions[j], orbs); aspect != nil {
				aspects = append(aspects, *aspect)
			}
		}
	}

	return aspects
}

// findAspect checks if two positions form an aspect
func findAspect(p1, p2 position.Position, orbs Orbs) *Aspect {
	// Calculate the angle between the two bodies
	angle := math.Abs(p1.EclipticLongitude - p2.EclipticLongitude)
	if angle > 180 {
		angle = 360 - angle
	}

	// Check each aspect type
	aspectTypes := []AspectType{Conjunction, Sextile, Square, Trine, Opposition}

	for _, aspectType := range aspectTypes {
		exactAngle := aspectType.Angle()
		orb := math.Abs(angle - exactAngle)

		maxOrb := orbs[aspectType]
		if maxOrb == 0 {
			maxOrb = DefaultOrbs[aspectType]
		}

		if orb <= maxOrb {
			return &Aspect{
				Body1: p1.Body,
				Body2: p2.Body,
				Type:  aspectType,
				Angle: angle,
				Orb:   orb,
			}
		}
	}

	return nil
}

// AspectBetween calculates the aspect (if any) between two specific bodies
func AspectBetween(p1, p2 position.Position, orbs Orbs) *Aspect {
	return findAspect(p1, p2, orbs)
}

// AllAspectTypes returns all aspect types
func AllAspectTypes() []AspectType {
	return []AspectType{Conjunction, Sextile, Square, Trine, Opposition}
}
