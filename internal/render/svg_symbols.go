package render

import (
	"embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"

	"github.com/ctrl-vfr/astral-tui/pkg/horoscope"
	"github.com/ctrl-vfr/astral-tui/pkg/position"
)

//go:embed symbols/zodiac/*.svg symbols/planets/*.svg
var symbolsFS embed.FS

// Font metrics from Noto Sans Symbols (unitsPerEm = 1000)
const glyphUnitsPerEm = 1000.0

var zodiacPaths map[horoscope.ZodiacSign]string
var planetPaths map[position.CelestialBody]string

// pathAttrRegex extracts the d attribute from SVG path element
var pathAttrRegex = regexp.MustCompile(`<path[^>]*\sd="([^"]*)"`)

func init() {
	zodiacPaths = make(map[horoscope.ZodiacSign]string)
	planetPaths = make(map[position.CelestialBody]string)

	zodiacFiles := map[string]horoscope.ZodiacSign{
		"aries.svg":       horoscope.Aries,
		"taurus.svg":      horoscope.Taurus,
		"gemini.svg":      horoscope.Gemini,
		"cancer.svg":      horoscope.Cancer,
		"leo.svg":         horoscope.Leo,
		"virgo.svg":       horoscope.Virgo,
		"libra.svg":       horoscope.Libra,
		"scorpio.svg":     horoscope.Scorpio,
		"sagittarius.svg": horoscope.Sagittarius,
		"capricorn.svg":   horoscope.Capricorn,
		"aquarius.svg":    horoscope.Aquarius,
		"pisces.svg":      horoscope.Pisces,
	}

	planetFiles := map[string]position.CelestialBody{
		"sun.svg":       position.Sun,
		"moon.svg":      position.Moon,
		"mercury.svg":   position.Mercury,
		"venus.svg":     position.Venus,
		"mars.svg":      position.Mars,
		"jupiter.svg":   position.Jupiter,
		"saturn.svg":    position.Saturn,
		"uranus.svg":    position.Uranus,
		"neptune.svg":   position.Neptune,
		"pluto.svg":     position.Pluto,
		"northnode.svg": position.NorthNode,
		"southnode.svg": position.SouthNode,
		"ceres.svg":     position.Ceres,
		"pallas.svg":    position.Pallas,
		"juno.svg":      position.Juno,
		"vesta.svg":     position.Vesta,
	}

	for filename, sign := range zodiacFiles {
		if path := loadPathFromSVG("symbols/zodiac/" + filename); path != "" {
			zodiacPaths[sign] = path
		}
	}

	for filename, body := range planetFiles {
		if path := loadPathFromSVG("symbols/planets/" + filename); path != "" {
			planetPaths[body] = path
		}
	}
}

func loadPathFromSVG(filepath string) string {
	data, err := symbolsFS.ReadFile(filepath)
	if err != nil {
		return ""
	}

	matches := pathAttrRegex.FindSubmatch(data)
	if len(matches) < 2 {
		return ""
	}

	return string(matches[1])
}

// pathCmdRegex matches SVG path commands with their coordinates
var pathCmdRegex = regexp.MustCompile(`([MLHVCSQTAZ])([^MLHVCSQTAZ]*)`)

// drawSymbol draws an SVG path symbol at the given position with scaling and color.
func drawSymbol(canvas *svg.SVG, path string, cx, cy int, size float64, color string) {
	if path == "" {
		return
	}

	scale := size / glyphUnitsPerEm
	transformedPath := transformPath(path, cx, cy, scale)
	canvas.Path(transformedPath, fmt.Sprintf("fill:%s", color))
}

// transformPath scales and translates an SVG path to center it at (cx, cy)
func transformPath(path string, cx, cy int, scale float64) string {
	glyphCenterX := 435.0
	glyphCenterY := 350.0

	var result strings.Builder
	matches := pathCmdRegex.FindAllStringSubmatch(path, -1)

	for _, match := range matches {
		cmd := match[1]
		coords := strings.TrimSpace(match[2])

		if coords == "" {
			result.WriteString(cmd)
			continue
		}

		nums := parseNumbers(coords)
		result.WriteString(cmd)

		switch cmd {
		case "M", "L", "T":
			for i := 0; i < len(nums); i += 2 {
				if i+1 < len(nums) {
					x := float64(cx) + (nums[i]-glyphCenterX)*scale
					y := float64(cy) - (nums[i+1]-glyphCenterY)*scale
					if i > 0 {
						result.WriteString(" ")
					}
					result.WriteString(fmt.Sprintf("%.1f %.1f", x, y))
				}
			}
		case "H":
			for i, n := range nums {
				x := float64(cx) + (n-glyphCenterX)*scale
				if i > 0 {
					result.WriteString(" ")
				}
				result.WriteString(fmt.Sprintf("%.1f", x))
			}
		case "V":
			for i, n := range nums {
				y := float64(cy) - (n-glyphCenterY)*scale
				if i > 0 {
					result.WriteString(" ")
				}
				result.WriteString(fmt.Sprintf("%.1f", y))
			}
		case "C":
			for i := 0; i < len(nums); i += 6 {
				if i+5 < len(nums) {
					x1 := float64(cx) + (nums[i]-glyphCenterX)*scale
					y1 := float64(cy) - (nums[i+1]-glyphCenterY)*scale
					x2 := float64(cx) + (nums[i+2]-glyphCenterX)*scale
					y2 := float64(cy) - (nums[i+3]-glyphCenterY)*scale
					x := float64(cx) + (nums[i+4]-glyphCenterX)*scale
					y := float64(cy) - (nums[i+5]-glyphCenterY)*scale
					if i > 0 {
						result.WriteString(" ")
					}
					result.WriteString(fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f %.1f", x1, y1, x2, y2, x, y))
				}
			}
		case "S":
			for i := 0; i < len(nums); i += 4 {
				if i+3 < len(nums) {
					x2 := float64(cx) + (nums[i]-glyphCenterX)*scale
					y2 := float64(cy) - (nums[i+1]-glyphCenterY)*scale
					x := float64(cx) + (nums[i+2]-glyphCenterX)*scale
					y := float64(cy) - (nums[i+3]-glyphCenterY)*scale
					if i > 0 {
						result.WriteString(" ")
					}
					result.WriteString(fmt.Sprintf("%.1f %.1f %.1f %.1f", x2, y2, x, y))
				}
			}
		case "Q":
			for i := 0; i < len(nums); i += 4 {
				if i+3 < len(nums) {
					x1 := float64(cx) + (nums[i]-glyphCenterX)*scale
					y1 := float64(cy) - (nums[i+1]-glyphCenterY)*scale
					x := float64(cx) + (nums[i+2]-glyphCenterX)*scale
					y := float64(cy) - (nums[i+3]-glyphCenterY)*scale
					if i > 0 {
						result.WriteString(" ")
					}
					result.WriteString(fmt.Sprintf("%.1f %.1f %.1f %.1f", x1, y1, x, y))
				}
			}
		case "A":
			for i := 0; i < len(nums); i += 7 {
				if i+6 < len(nums) {
					rx := nums[i] * scale
					ry := nums[i+1] * scale
					rotation := nums[i+2]
					largeArc := nums[i+3]
					sweep := nums[i+4]
					x := float64(cx) + (nums[i+5]-glyphCenterX)*scale
					y := float64(cy) - (nums[i+6]-glyphCenterY)*scale
					if i > 0 {
						result.WriteString(" ")
					}
					result.WriteString(fmt.Sprintf("%.1f %.1f %.0f %.0f %.0f %.1f %.1f", rx, ry, rotation, largeArc, sweep, x, y))
				}
			}
		case "Z":
			// Close path - no coordinates
		}
	}

	return result.String()
}

func parseNumbers(s string) []float64 {
	var nums []float64
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == ','
	})

	for _, p := range parts {
		subParts := splitOnNegative(p)
		for _, sp := range subParts {
			if sp != "" {
				if n, err := strconv.ParseFloat(sp, 64); err == nil {
					nums = append(nums, n)
				}
			}
		}
	}
	return nums
}

func splitOnNegative(s string) []string {
	var result []string
	var current strings.Builder

	for i, r := range s {
		if r == '-' && i > 0 && s[i-1] != 'E' && s[i-1] != 'e' {
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		}
		current.WriteRune(r)
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// GetZodiacPath returns the SVG path for a zodiac sign
func GetZodiacPath(sign horoscope.ZodiacSign) string {
	return zodiacPaths[sign]
}

// GetPlanetPath returns the SVG path for a celestial body
func GetPlanetPath(body position.CelestialBody) string {
	return planetPaths[body]
}
