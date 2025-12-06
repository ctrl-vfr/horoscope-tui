package horoscope

import (
	"fmt"
)

// ZodiacSign represents one of the 12 zodiac signs
type ZodiacSign int

// The twelve zodiac signs in order.
const (
	Aries ZodiacSign = iota
	Taurus
	Gemini
	Cancer
	Leo
	Virgo
	Libra
	Scorpio
	Sagittarius
	Capricorn
	Aquarius
	Pisces
)

// String returns the name of the zodiac sign
func (z ZodiacSign) String() string {
	return signNames[z]
}

// Symbol returns the astrological symbol for the sign
func (z ZodiacSign) Symbol() string {
	return signSymbols[z]
}

// Element returns the element (Fire, Earth, Air, Water) for the sign
func (z ZodiacSign) Element() Element {
	return signElements[z]
}

// Modality returns the modality (Cardinal, Fixed, Mutable) for the sign
func (z ZodiacSign) Modality() Modality {
	return signModalities[z]
}

var signNames = map[ZodiacSign]string{
	Aries:       "Aries",
	Taurus:      "Taurus",
	Gemini:      "Gemini",
	Cancer:      "Cancer",
	Leo:         "Leo",
	Virgo:       "Virgo",
	Libra:       "Libra",
	Scorpio:     "Scorpio",
	Sagittarius: "Sagittarius",
	Capricorn:   "Capricorn",
	Aquarius:    "Aquarius",
	Pisces:      "Pisces",
}

var signSymbols = map[ZodiacSign]string{
	Aries:       "♈",
	Taurus:      "♉",
	Gemini:      "♊",
	Cancer:      "♋",
	Leo:         "♌",
	Virgo:       "♍",
	Libra:       "♎",
	Scorpio:     "♏",
	Sagittarius: "♐",
	Capricorn:   "♑",
	Aquarius:    "♒",
	Pisces:      "♓",
}

// Element represents the four classical elements
type Element int

// The four classical elements.
const (
	Fire Element = iota
	Earth
	Air
	Water
)

// String returns the name of the element.
func (e Element) String() string {
	return []string{"Fire", "Earth", "Air", "Water"}[e]
}

var signElements = map[ZodiacSign]Element{
	Aries: Fire, Leo: Fire, Sagittarius: Fire,
	Taurus: Earth, Virgo: Earth, Capricorn: Earth,
	Gemini: Air, Libra: Air, Aquarius: Air,
	Cancer: Water, Scorpio: Water, Pisces: Water,
}

// Modality represents the three modalities
type Modality int

// The three modalities.
const (
	Cardinal Modality = iota
	Fixed
	Mutable
)

// String returns the name of the modality.
func (m Modality) String() string {
	return []string{"Cardinal", "Fixed", "Mutable"}[m]
}

var signModalities = map[ZodiacSign]Modality{
	Aries: Cardinal, Cancer: Cardinal, Libra: Cardinal, Capricorn: Cardinal,
	Taurus: Fixed, Leo: Fixed, Scorpio: Fixed, Aquarius: Fixed,
	Gemini: Mutable, Virgo: Mutable, Sagittarius: Mutable, Pisces: Mutable,
}

// ZodiacPosition represents a position within a zodiac sign
type ZodiacPosition struct {
	Sign    ZodiacSign
	Degrees int     // 0-29
	Minutes int     // 0-59
	Seconds int     // 0-59
	Total   float64 // Total degrees in sign (0-30)
}

// String returns a formatted position string (e.g., "15°23'45" Aries")
func (z ZodiacPosition) String() string {
	return fmt.Sprintf("%d°%02d' %s", z.Degrees, z.Minutes, z.Sign.String())
}

// ShortString returns abbreviated position (e.g., "15°23' Ari")
func (z ZodiacPosition) ShortString() string {
	abbrev := z.Sign.String()[:3]
	return fmt.Sprintf("%d°%02d' %s", z.Degrees, z.Minutes, abbrev)
}

// LongitudeToZodiac converts ecliptic longitude (0-360) to zodiac position
func LongitudeToZodiac(longitude float64) ZodiacPosition {
	// Normalize to 0-360
	for longitude < 0 {
		longitude += 360
	}
	for longitude >= 360 {
		longitude -= 360
	}

	// Determine sign (each sign is 30 degrees)
	signIndex := int(longitude / 30)
	if signIndex > 11 {
		signIndex = 11
	}

	// Position within sign
	posInSign := longitude - float64(signIndex*30)

	// Convert to degrees, minutes, seconds
	degrees := int(posInSign)
	fracDeg := posInSign - float64(degrees)
	minutes := int(fracDeg * 60)
	fracMin := fracDeg*60 - float64(minutes)
	seconds := int(fracMin * 60)

	return ZodiacPosition{
		Sign:    ZodiacSign(signIndex),
		Degrees: degrees,
		Minutes: minutes,
		Seconds: seconds,
		Total:   posInSign,
	}
}

// AllSigns returns all zodiac signs in order
func AllSigns() []ZodiacSign {
	return []ZodiacSign{
		Aries, Taurus, Gemini, Cancer, Leo, Virgo,
		Libra, Scorpio, Sagittarius, Capricorn, Aquarius, Pisces,
	}
}
