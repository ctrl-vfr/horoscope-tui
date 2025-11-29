package render

import (
	"bytes"

	svg "github.com/ajstarks/svgo"

	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

// SVGWheelGenerator generates SVG zodiac wheels.
type SVGWheelGenerator struct {
	size   int
	center int
	radius int
}

// NewSVGWheelGenerator creates a new SVG generator.
func NewSVGWheelGenerator(size int) *SVGWheelGenerator {
	return &SVGWheelGenerator{
		size:   size,
		center: size / 2,
		radius: size/2 - 20,
	}
}

// Generate creates an SVG zodiac wheel (natal only).
func (g *SVGWheelGenerator) Generate(positions []position.Position) []byte {
	return g.GenerateWithTransits(positions, nil)
}

// GenerateWithTransits creates an SVG zodiac wheel with natal and transit positions.
func (g *SVGWheelGenerator) GenerateWithTransits(natal, transits []position.Position) []byte {
	var buf bytes.Buffer
	canvas := svg.New(&buf)

	canvas.Start(g.size, g.size)
	g.drawOuterWheel(canvas)
	g.drawZodiacSegments(canvas)
	g.drawInnerCircle(canvas)
	g.drawAxes(canvas)
	g.drawPlanetsRing(canvas, natal, 0.40, false)
	if len(transits) > 0 {
		g.drawPlanetsRing(canvas, transits, 0.55, true)
	}
	g.drawAspects(canvas, natal)
	canvas.End()

	return buf.Bytes()
}
