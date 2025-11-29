package messages

import (
	"time"

	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
)

// Form messages
type FormSubmittedMsg struct {
	Date time.Time
	City string
}

type DateChangedMsg struct {
	Date time.Time
}

// Geocoding messages
type GeocodingResultMsg struct {
	Latitude    float64
	Longitude   float64
	DisplayName string
	Err         error
}

// Chart calculation messages
type ChartReadyMsg struct {
	Chart *horoscope.Chart
}

type ChartErrorMsg struct {
	Err error
}

// Wheel generation messages
type WheelGeneratedMsg struct {
	PNGData []byte
	Err     error
}

// GPT interpretation messages
type InterpReadyMsg struct {
	Content string
	Err     error
}

// Status messages
type StatusMsg string

type ErrorMsg struct {
	Err error
}
