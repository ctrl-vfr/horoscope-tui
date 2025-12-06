package render

import (
	"fmt"
	"math"
	"sort"

	svg "github.com/ajstarks/svgo"

	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
	"github.com/ctrl-vfr/horoscope-tui/pkg/position"
)

func (g *SVGWheelGenerator) drawOuterWheel(canvas *svg.SVG) {
	canvas.Circle(g.center, g.center, g.radius, fmt.Sprintf("fill:none;stroke:%s;stroke-width:1;opacity:0.3", svgPrimary))
	canvas.Circle(g.center, g.center, g.radius, fmt.Sprintf("fill:none;stroke:%s;stroke-width:1", svgPrimary))
	canvas.Circle(g.center, g.center, g.radius-54, fmt.Sprintf("fill:none;stroke:%s;stroke-width:1", svgBorder))
}

func (g *SVGWheelGenerator) drawZodiacSegments(canvas *svg.SVG) {
	signs := []struct {
		sign  horoscope.ZodiacSign
		start float64
	}{
		{horoscope.Aries, 0},
		{horoscope.Taurus, 30},
		{horoscope.Gemini, 60},
		{horoscope.Cancer, 90},
		{horoscope.Leo, 120},
		{horoscope.Virgo, 150},
		{horoscope.Libra, 180},
		{horoscope.Scorpio, 210},
		{horoscope.Sagittarius, 240},
		{horoscope.Capricorn, 270},
		{horoscope.Aquarius, 300},
		{horoscope.Pisces, 330},
	}

	for _, s := range signs {
		angle := (90 - s.start) * math.Pi / 180
		x1 := g.center + int(float64(g.radius-54)*math.Cos(angle))
		y1 := g.center - int(float64(g.radius-54)*math.Sin(angle))
		x2 := g.center + int(float64(g.radius)*math.Cos(angle))
		y2 := g.center - int(float64(g.radius)*math.Sin(angle))
		canvas.Line(x1, y1, x2, y2, fmt.Sprintf("stroke:%s;stroke-width:2", svgBorder))

		midAngle := (90 - s.start - 15) * math.Pi / 180
		symbolRadius := float64(g.radius) - 25
		sx := g.center + int(symbolRadius*math.Cos(midAngle))
		sy := g.center - int(symbolRadius*math.Sin(midAngle))

		color := getElementColor(s.sign.Element())
		drawSymbol(canvas, GetZodiacPath(s.sign), sx, sy+7, 34, color)
	}
}

func (g *SVGWheelGenerator) drawInnerCircle(canvas *svg.SVG) {
	innerRadius := int(float64(g.radius) * 0.65)
	canvas.Circle(g.center, g.center, innerRadius, fmt.Sprintf("fill:none;stroke:%s;stroke-width:1", svgPrimary))
	canvas.Circle(g.center, g.center, 8, fmt.Sprintf("fill:%s;opacity:0.3", svgBright))
	canvas.Circle(g.center, g.center, 4, fmt.Sprintf("fill:%s", svgBright))
}

func (g *SVGWheelGenerator) drawAxes(canvas *svg.SVG) {
	innerRadius := int(float64(g.radius) * 0.65)
	canvas.Line(g.center, g.center-innerRadius, g.center, g.center+innerRadius, fmt.Sprintf("stroke:%s;stroke-width:1;stroke-dasharray:5,3", svgPurple))
	canvas.Line(g.center-innerRadius, g.center, g.center+innerRadius, g.center, fmt.Sprintf("stroke:%s;stroke-width:1;stroke-dasharray:5,3", svgPurple))
	canvas.Text(g.center, g.center-innerRadius-8, "MC", fmt.Sprintf("font-size:10px;fill:%s;text-anchor:middle", svgTextLight))
	canvas.Text(g.center, g.center+innerRadius+14, "IC", fmt.Sprintf("font-size:10px;fill:%s;text-anchor:middle", svgTextLight))
	canvas.Text(g.center+innerRadius+10, g.center+4, "ASC", fmt.Sprintf("font-size:10px;fill:%s;text-anchor:start", svgTextLight))
	canvas.Text(g.center-innerRadius-10, g.center+4, "DSC", fmt.Sprintf("font-size:10px;fill:%s;text-anchor:end", svgTextLight))
}

func (g *SVGWheelGenerator) drawPlanetsRing(canvas *svg.SVG, positions []position.Position, radiusFactor float64, isTransit bool) {
	planetRadius := float64(g.radius) * radiusFactor

	maxBody := position.Pluto
	if isTransit {
		maxBody = position.Neptune
	}

	// Filter and sort positions by longitude
	var filtered []position.Position
	for _, pos := range positions {
		if pos.Body <= maxBody {
			filtered = append(filtered, pos)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].EclipticLongitude < filtered[j].EclipticLongitude
	})

	// Calculate radial offsets to avoid overlapping (threshold: 5°)
	offsets := calculateRadialOffsets(filtered, 15.0)

	for i, pos := range filtered {
		angle := (90 - pos.EclipticLongitude) * math.Pi / 180
		adjustedRadius := planetRadius + offsets[i]

		x := g.center + int(adjustedRadius*math.Cos(angle))
		y := g.center - int(adjustedRadius*math.Sin(angle))

		var color string
		var size float64
		if isTransit {
			color = svgAccent
			size = 20.0
		} else {
			color = getPlanetSVGColor(pos.Body)
			size = 28.0
		}
		drawSymbol(canvas, GetPlanetPath(pos.Body), x, y, size, color)
	}
}

func calculateRadialOffsets(positions []position.Position, threshold float64) []float64 {
	offsets := make([]float64, len(positions))
	const offsetStep = -25.0

	for i := 1; i < len(positions); i++ {
		diff := positions[i].EclipticLongitude - positions[i-1].EclipticLongitude
		if diff < 0 {
			diff += 360
		}
		if diff > 180 {
			diff = 360 - diff
		}

		if diff < threshold {
			offsets[i] = offsets[i-1] + offsetStep
			if offsets[i] < -60 {
				offsets[i] = -60
			}
		}
	}
	return offsets
}

func (g *SVGWheelGenerator) drawAspects(canvas *svg.SVG, positions []position.Position) {
	aspects := horoscope.CalculateAspects(positions, horoscope.TightOrbs)
	planetRadius := float64(g.radius) * 0.35

	for _, aspect := range aspects {
		if aspect.Body1 > position.Saturn || aspect.Body2 > position.Saturn {
			continue
		}

		var lon1, lon2 float64
		for _, pos := range positions {
			if pos.Body == aspect.Body1 {
				lon1 = pos.EclipticLongitude
			}
			if pos.Body == aspect.Body2 {
				lon2 = pos.EclipticLongitude
			}
		}

		angle1 := (90 - lon1) * math.Pi / 180
		angle2 := (90 - lon2) * math.Pi / 180

		x1 := g.center + int(planetRadius*math.Cos(angle1))
		y1 := g.center - int(planetRadius*math.Sin(angle1))
		x2 := g.center + int(planetRadius*math.Cos(angle2))
		y2 := g.center - int(planetRadius*math.Sin(angle2))

		color := getAspectColor(aspect.Type)
		opacity := 0.6 - aspect.Orb*0.05
		if opacity < 0.2 {
			opacity = 0.2
		}

		style := fmt.Sprintf("stroke:%s;stroke-width:1;opacity:%.1f", color, opacity)
		if aspect.Type == horoscope.Trine || aspect.Type == horoscope.Sextile {
			style += ";stroke-dasharray:4,2"
		}

		canvas.Line(x1, y1, x2, y2, style)
	}
}
