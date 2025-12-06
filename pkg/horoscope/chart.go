// Package horoscope provides astrological calculations
package horoscope

import (
	"time"

	"github.com/ctrl-vfr/astral-tui/pkg/position"
)

// Chart represents a complete astrological chart
type Chart struct {
	DateTime  time.Time
	Latitude  float64
	Longitude float64
	Location  string

	Positions []position.Position
	Houses    HouseCusps
	Aspects   []Aspect
}

// HouseCusps interface for house calculation results
type HouseCusps interface {
	GetHouse(longitude float64) int
}

// BodyInHouse returns the house number for a given body
func (c *Chart) BodyInHouse(body position.CelestialBody) int {
	for _, pos := range c.Positions {
		if pos.Body == body {
			if c.Houses != nil {
				return c.Houses.GetHouse(pos.EclipticLongitude)
			}
			return 0
		}
	}
	return 0
}

// GetPosition returns the position of a specific body
func (c *Chart) GetPosition(body position.CelestialBody) *position.Position {
	for _, pos := range c.Positions {
		if pos.Body == body {
			return &pos
		}
	}
	return nil
}

// GetZodiacPosition returns the zodiac position for a body
func (c *Chart) GetZodiacPosition(body position.CelestialBody) *ZodiacPosition {
	pos := c.GetPosition(body)
	if pos == nil {
		return nil
	}
	zp := LongitudeToZodiac(pos.EclipticLongitude)
	return &zp
}

// AspectsFor returns all aspects involving a specific body
func (c *Chart) AspectsFor(body position.CelestialBody) []Aspect {
	var result []Aspect
	for _, aspect := range c.Aspects {
		if aspect.Body1 == body || aspect.Body2 == body {
			result = append(result, aspect)
		}
	}
	return result
}
