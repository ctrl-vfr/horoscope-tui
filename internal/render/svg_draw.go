package render

import (
	"fmt"
	"math"

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
		canvas.Text(sx, sy+14, s.sign.Symbol(), fmt.Sprintf("font-size:34px;fill:%s;text-anchor:middle;font-family:sans-serif", color))
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
	placed := make(map[int]bool)

	for _, pos := range positions {
		maxBody := position.Pluto
		if isTransit {
			maxBody = position.Neptune
		}
		if pos.Body > maxBody {
			continue
		}

		angle := (90 - pos.EclipticLongitude) * math.Pi / 180

		offset := 0.0
		for {
			x := g.center + int((planetRadius+offset)*math.Cos(angle))
			y := g.center - int((planetRadius+offset)*math.Sin(angle))
			key := x*1000 + y
			if !placed[key] {
				placed[key] = true

				var color, fontSize string
				if isTransit {
					color = svgSecondary
					fontSize = "20px"
				} else {
					color = getPlanetSVGColor(pos.Body)
					fontSize = "30px"
				}

				canvas.Text(x, y+5, pos.Body.Symbol(), fmt.Sprintf("font-size:%s;fill:%s;text-anchor:middle;font-family:sans-serif", fontSize, color))
				break
			}
			offset -= 12
			if offset < -48 {
				break
			}
		}
	}
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
