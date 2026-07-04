package format

import (
	"fmt"
	"strings"
	"time"
)

// Common format presets
var presets = map[string]string{
	"rfc3339":       time.RFC3339,
	"rfc3339nano":   time.RFC3339Nano,
	"rfc1123":       time.RFC1123,
	"rfc1123z":      time.RFC1123Z,
	"rfc822":        time.RFC822,
	"rfc822z":       time.RFC822Z,
	"kitchen":       time.Kitchen,
	"stamp":         time.Stamp,
	"stampmilli":    time.StampMilli,
	"stampmicro":    time.StampMicro,
	"stampnano":     time.StampNano,
	"datetime":      "2006-01-02 15:04:05",
	"date":          "2006-01-02",
	"time":          "15:04:05",
	"time-ms":       "15:04:05.000",
	"unix":          "unix",
	"unix-ms":       "unix-ms",
	"unix-ns":       "unix-ns",
	"iso8601":       "2006-01-02T15:04:05Z07:00",
	"log":           "2006/01/02 15:04:05",
	"log-ms":        "2006/01/02 15:04:05.000",
	"sql":           "2006-01-02 15:04:05",
	"filename":      "20060102-150405",
	"filename-ms":   "20060102-150405.000",
	"shortdate":     "01/02/2006",
	"us-shortdate":  "01/02/2006",
	"eu-shortdate":  "02/01/2006",
	"year-month":    "2006-01",
	"month-day":     "01-02",
	"month":         "January",
	"month-short":   "Jan",
	"weekday":       "Monday",
	"weekday-short": "Mon",
	"timezone":      "MST",
	"utc":           "UTC",
	"relative":      "relative",
}

// Format formats a time using a named preset or custom layout
func Format(t time.Time, layout string) (string, error) {
	layout = strings.TrimSpace(strings.ToLower(layout))

	// Handle special formats
	switch layout {
	case "unix":
		return fmt.Sprintf("%d", t.Unix()), nil
	case "unix-ms":
		return fmt.Sprintf("%d", t.UnixMilli()), nil
	case "unix-ns":
		return fmt.Sprintf("%d", t.UnixNano()), nil
	case "relative":
		return FormatRelative(t), nil
	}

	// Check presets
	if fmt, ok := presets[layout]; ok {
		return t.Format(fmt), nil
	}

	// Try as Go layout directly
	return t.Format(layout), nil
}

// FormatRelative formats a time as a relative string ("2 hours ago", "in 3 days")
func FormatRelative(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	abs := diff
	if abs < 0 {
		abs = -abs
	}

	switch {
	case abs < time.Second:
		return "just now"
	case abs < time.Minute:
		secs := int(abs.Seconds())
		if diff > 0 {
			return fmt.Sprintf("%d second%s ago", secs, plural(secs))
		}
		return fmt.Sprintf("in %d second%s", secs, plural(secs))
	case abs < time.Hour:
		mins := int(abs.Minutes())
		if diff > 0 {
			return fmt.Sprintf("%d minute%s ago", mins, plural(mins))
		}
		return fmt.Sprintf("in %d minute%s", mins, plural(mins))
	case abs < 24*time.Hour:
		hours := int(abs.Hours())
		if diff > 0 {
			return fmt.Sprintf("%d hour%s ago", hours, plural(hours))
		}
		return fmt.Sprintf("in %d hour%s", hours, plural(hours))
	case abs < 30*24*time.Hour:
		days := int(abs.Hours() / 24)
		if diff > 0 {
			return fmt.Sprintf("%d day%s ago", days, plural(days))
		}
		return fmt.Sprintf("in %d day%s", days, plural(days))
	case abs < 365*24*time.Hour:
		months := int(abs.Hours() / (24 * 30))
		if diff > 0 {
			return fmt.Sprintf("%d month%s ago", months, plural(months))
		}
		return fmt.Sprintf("in %d month%s", months, plural(months))
	default:
		years := int(abs.Hours() / (24 * 365))
		if diff > 0 {
			return fmt.Sprintf("%d year%s ago", years, plural(years))
		}
		return fmt.Sprintf("in %d year%s", years, plural(years))
	}
}

// ListFormats returns all available format presets
func ListFormats() map[string]string {
	result := make(map[string]string)
	for k, v := range presets {
		result[k] = v
	}
	return result
}

// GetPreset returns the Go layout for a preset name
func GetPreset(name string) (string, bool) {
	layout, ok := presets[strings.ToLower(name)]
	return layout, ok
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
