// Package i18n provides internationalization support with automatic locale detection.
package i18n

import (
	"os"
	"strings"
)

// Lang represents a supported language.
type Lang string

// Supported languages.
const (
	EN Lang = "en"
	FR Lang = "fr"
	ES Lang = "es"
	DE Lang = "de"
)

var currentLang Lang

// Init detects the system locale and sets the current language.
// It checks LC_ALL, LC_MESSAGES, then LANG environment variables.
// Falls back to English if the locale is not supported.
func Init() {
	locale := os.Getenv("LC_ALL")
	if locale == "" {
		locale = os.Getenv("LC_MESSAGES")
	}
	if locale == "" {
		locale = os.Getenv("LANG")
	}

	currentLang = parseLocale(locale)
}

// parseLocale extracts the language code from a locale string.
// Examples: "fr_FR.UTF-8" -> FR, "en_US" -> EN, "de" -> DE
func parseLocale(locale string) Lang {
	if locale == "" {
		return EN
	}

	// Remove encoding suffix (e.g., ".UTF-8")
	if idx := strings.Index(locale, "."); idx != -1 {
		locale = locale[:idx]
	}

	// Extract language code (first 2 characters)
	lang := strings.ToLower(locale)
	if len(lang) >= 2 {
		lang = lang[:2]
	}

	switch lang {
	case "fr":
		return FR
	case "es":
		return ES
	case "de":
		return DE
	default:
		return EN
	}
}

// Current returns the currently active language.
func Current() Lang {
	return currentLang
}

// T returns the translation for the given key in the current language.
func T(key string) string {
	if msgs, ok := messages[currentLang]; ok {
		if val, ok := msgs[key]; ok {
			return val
		}
	}
	// Fallback to English
	if msgs, ok := messages[EN]; ok {
		if val, ok := msgs[key]; ok {
			return val
		}
	}
	return key
}

// SystemPrompt returns the localized system prompt for the OpenAI API.
func SystemPrompt() string {
	return systemPrompts[currentLang]
}

// Weekday returns the localized weekday name.
func Weekday(day int) string {
	key := weekdayKeys[day]
	return T(key)
}

var weekdayKeys = []string{
	"WeekdaySunday",
	"WeekdayMonday",
	"WeekdayTuesday",
	"WeekdayWednesday",
	"WeekdayThursday",
	"WeekdayFriday",
	"WeekdaySaturday",
}
