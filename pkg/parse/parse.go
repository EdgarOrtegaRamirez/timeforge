package parse

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Common timestamp formats, ordered by specificity
var formats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.000Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02T15:04:05Z",
	"2006-01-02",
	"01/02/2006 15:04:05",
	"01/02/2006",
	"02 Jan 2006 15:04:05",
	"02 Jan 2006",
	"Jan 2, 2006 15:04:05",
	"Jan 2, 2006",
	"Monday, 02-Jan-06 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"2006/01/02 15:04:05",
	"2006/01/02",
	"02/01/2006 15:04:05",
	"02/01/2006",
}

// Parse attempts to parse a timestamp string into a time.Time
// It tries multiple formats and also handles Unix timestamps (seconds, milliseconds, nanoseconds)
func Parse(input string) (time.Time, string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return time.Time{}, "", fmt.Errorf("empty input")
	}

	// Try Unix timestamps (numeric only)
	if t, err := parseUnixTimestamp(input); err == nil {
		return t, "unix", nil
	}

	// Try Go duration-style relative time
	if t, err := parseRelative(input); err == nil {
		return t, "relative", nil
	}

	// Try known formats
	for _, format := range formats {
		if t, err := time.Parse(format, input); err == nil {
			return t, format, nil
		}
	}

	return time.Time{}, "", fmt.Errorf("unable to parse timestamp: %q", input)
}

// ParseWithLocation parses a timestamp with a specific timezone
func ParseWithLocation(input, tz string) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, fmt.Errorf("unknown timezone %q: %w", tz, err)
	}

	t, _, err := Parse(input)
	if err != nil {
		return time.Time{}, err
	}

	return t.In(loc), nil
}

// parseUnixTimestamp tries to parse as a Unix timestamp (seconds, milliseconds, or nanoseconds)
func parseUnixTimestamp(input string) (time.Time, error) {
	// Try integer
	if v, err := strconv.ParseInt(input, 10, 64); err == nil {
		return unixToTime(v)
	}

	// Try float (for fractional seconds)
	if v, err := strconv.ParseFloat(input, 64); err == nil {
		sec := int64(v)
		nsec := int64((v - float64(sec)) * 1e9)
		return time.Unix(sec, nsec), nil
	}

	return time.Time{}, fmt.Errorf("not a unix timestamp")
}

func unixToTime(v int64) (time.Time, error) {
	switch {
	case v > 1e18: // nanoseconds
		return time.Unix(0, v), nil
	case v > 1e12: // milliseconds
		return time.UnixMilli(v), nil
	case v > 1e9: // seconds (beyond year 2001)
		return time.Unix(v, 0), nil
	case v > 0: // seconds (older dates)
		return time.Unix(v, 0), nil
	default:
		return time.Time{}, fmt.Errorf("invalid unix timestamp: %d", v)
	}
}

// parseRelative handles relative time expressions like "2h30m", "-1d", "+3h"
func parseRelative(input string) (time.Time, error) {
	// Handle signed durations like -1h30m or +2d
	signed := false
	s := input
	if strings.HasPrefix(s, "+") {
		s = s[1:]
		signed = true
	} else if strings.HasPrefix(s, "-") {
		s = s[1:]
		signed = true
	}

	if !signed {
		return time.Time{}, fmt.Errorf("not a relative time")
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return time.Time{}, err
	}

	if strings.HasPrefix(input, "-") {
		return time.Now().Add(-d), nil
	}
	return time.Now().Add(d), nil
}

// DetectFormat returns the detected format name for a timestamp string
func DetectFormat(input string) string {
	_, format, _ := Parse(input)
	return format
}

// IsTimestamp returns true if the input looks like a timestamp
func IsTimestamp(input string) bool {
	_, _, err := Parse(input)
	return err == nil
}
