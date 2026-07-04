package duration

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// DurationResult represents a computed duration with breakdown
type DurationResult struct {
	Duration  time.Duration
	Seconds   float64
	Minutes   float64
	Hours     float64
	Days      float64
	Weeks     float64
	Breakdown Breakdown
}

// Breakdown is a human-readable breakdown of a duration
type Breakdown struct {
	Days    int
	Hours   int
	Minutes int
	Seconds int
	Millis  int
}

// Compute calculates the duration between two timestamps
func Compute(start, end time.Time) *DurationResult {
	d := end.Sub(start)
	return FromDuration(d)
}

// FromDuration creates a DurationResult from a time.Duration
func FromDuration(d time.Duration) *DurationResult {
	secs := d.Seconds()
	mins := d.Minutes()
	hours := d.Hours()
	days := hours / 24
	weeks := days / 7

	// Break down into components
	abs := d
	if abs < 0 {
		abs = -abs
	}

	daysI := int(abs.Hours()) / 24
	hoursI := (int(abs.Hours())) % 24
	minsI := int(abs.Minutes()) % 60
	secsI := int(abs.Seconds()) % 60
	millisI := int(abs.Milliseconds()) % 1000

	return &DurationResult{
		Duration: d,
		Seconds:  secs,
		Minutes:  mins,
		Hours:    hours,
		Days:     days,
		Weeks:    weeks,
		Breakdown: Breakdown{
			Days:    daysI,
			Hours:   hoursI,
			Minutes: minsI,
			Seconds: secsI,
			Millis:  millisI,
		},
	}
}

// Add adds durations together
func Add(durations ...time.Duration) time.Duration {
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total
}

// Subtract subtracts durations
func Subtract(base time.Duration, subtract ...time.Duration) time.Duration {
	for _, d := range subtract {
		base -= d
	}
	return base
}

// Average calculates the average of multiple durations
func Average(durations ...time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

// Min returns the smallest duration
func Min(durations ...time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	m := durations[0]
	for _, d := range durations[1:] {
		if d < m {
			m = d
		}
	}
	return m
}

// Max returns the largest duration
func Max(durations ...time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	m := durations[0]
	for _, d := range durations[1:] {
		if d > m {
			m = d
		}
	}
	return m
}

// FormatCompact formats a duration compactly: "2h30m15s"
func FormatCompact(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	abs := d
	if abs < 0 {
		abs = -abs
	}

	var parts []string

	days := int(abs.Hours()) / 24
	hours := int(abs.Hours()) % 24
	mins := int(abs.Minutes()) % 60
	secs := int(abs.Seconds()) % 60
	millis := int(abs.Milliseconds()) % 1000

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if mins > 0 {
		parts = append(parts, fmt.Sprintf("%dm", mins))
	}
	if secs > 0 {
		parts = append(parts, fmt.Sprintf("%ds", secs))
	}
	if millis > 0 && len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%dms", millis))
	}

	result := strings.Join(parts, "")
	if d < 0 {
		result = "-" + result
	}
	return result
}

// FormatHuman formats a duration in human-readable form
func FormatHuman(d time.Duration) string {
	if d == 0 {
		return "zero"
	}

	abs := d
	if abs < 0 {
		abs = -abs
	}

	days := int(abs.Hours()) / 24
	hours := int(abs.Hours()) % 24
	mins := int(abs.Minutes()) % 60
	secs := int(abs.Seconds()) % 60

	var parts []string

	if days > 365 {
		years := days / 365
		days = days % 365
		parts = append(parts, fmt.Sprintf("%d year%s", years, plural(years)))
	}
	if days > 30 {
		months := days / 30
		days = days % 30
		parts = append(parts, fmt.Sprintf("%d month%s", months, plural(months)))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d day%s", days, plural(days)))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hour%s", hours, plural(hours)))
	}
	if mins > 0 {
		parts = append(parts, fmt.Sprintf("%d minute%s", mins, plural(mins)))
	}
	if secs > 0 {
		parts = append(parts, fmt.Sprintf("%d second%s", secs, plural(secs)))
	}

	if len(parts) == 0 {
		return "less than a second"
	}

	result := strings.Join(parts, ", ")
	if d < 0 {
		result = "-" + result
	}
	return result
}

// Parse parses a duration string with extended syntax (supports days/weeks/months)
func Parse(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty duration")
	}

	// Try standard Go duration first
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	// Extended format: parse days, weeks, months manually
	var total time.Duration
	remaining := s

	for remaining != "" {
		var value int
		var unit string
		n, err := fmt.Sscanf(remaining, "%d%s", &value, &unit)
		if n < 2 || err != nil {
			return 0, fmt.Errorf("invalid duration part in %q", s)
		}

		// Find where this part ends
		idx := strings.IndexAny(remaining, "0123456789")
		if idx >= 0 {
			// Find end of number
			numEnd := idx
			for numEnd < len(remaining) && remaining[numEnd] >= '0' && remaining[numEnd] <= '9' {
				numEnd++
			}
			// Find end of unit
			unitEnd := numEnd
			for unitEnd < len(remaining) && remaining[unitEnd] != ' ' && !(remaining[unitEnd] >= '0' && remaining[unitEnd] <= '9') {
				unitEnd++
			}
			remaining = remaining[unitEnd:]
		} else {
			break
		}

		switch unit {
		case "y", "yr", "yrs", "year", "years":
			total += time.Duration(value) * 365 * 24 * time.Hour
		case "M", "mo", "mos", "month", "months":
			total += time.Duration(value) * 30 * 24 * time.Hour
		case "w", "wk", "wks", "week", "weeks":
			total += time.Duration(value) * 7 * 24 * time.Hour
		case "d", "day", "days":
			total += time.Duration(value) * 24 * time.Hour
		case "h", "hr", "hrs", "hour", "hours":
			total += time.Duration(value) * time.Hour
		case "m", "min", "mins", "minute", "minutes":
			total += time.Duration(value) * time.Minute
		case "s", "sec", "secs", "second", "seconds":
			total += time.Duration(value) * time.Second
		case "ms", "milli", "millis":
			total += time.Duration(value) * time.Millisecond
		default:
			return 0, fmt.Errorf("unknown duration unit: %q", unit)
		}
	}

	return total, nil
}

// Difference returns the absolute difference between two durations
func Difference(a, b time.Duration) time.Duration {
	d := a - b
	if d < 0 {
		return -d
	}
	return d
}

// Percentage calculates what percentage b is of a
func Percentage(a, b time.Duration) float64 {
	if a == 0 {
		return math.Inf(1)
	}
	return (float64(b) / float64(a)) * 100
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
