// Package messages defines the message types used for TUI component communication.
package messages

import (
	"time"

	"github.com/ctrl-vfr/horoscope-tui/pkg/horoscope"
)

// FormSubmittedMsg is sent when the user submits the form
type FormSubmittedMsg struct {
	Date time.Time
	City string
}

// DateChangedMsg is sent when the user changes the date
type DateChangedMsg struct {
	Date time.Time
}

// TransitDateChangedMsg is sent when the user changes the transit date
type TransitDateChangedMsg struct {
	Date time.Time
}

// GeocodingResultMsg is sent when the user submits the form
type GeocodingResultMsg struct {
	Latitude    float64
	Longitude   float64
	DisplayName string
	Err         error
}

// ChartReadyMsg is sent when the chart is ready
type ChartReadyMsg struct {
	Chart *horoscope.Chart
}

// ChartErrorMsg is sent when the chart generation fails
type ChartErrorMsg struct {
	Err error
}

// WheelGeneratedMsg is sent when the wheel is generated
type WheelGeneratedMsg struct {
	PNGData []byte
	Err     error
}

// InterpReadyMsg is sent when the interpretation is ready
type InterpReadyMsg struct {
	Content string
	Err     error
}

// StatusMsg is sent when the status changes
type StatusMsg string

// ErrorMsg is sent when an error occurs
type ErrorMsg struct {
	Err error
}
