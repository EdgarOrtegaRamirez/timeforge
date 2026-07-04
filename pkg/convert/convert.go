package convert

import (
	"fmt"
	"regexp"
	"time"
)

// TimezoneInfo contains information about a timezone
type TimezoneInfo struct {
	Name   string
	Offset int // UTC offset in seconds
	Abbrev string
	IsDST  bool
}

// ConvertTimezone converts a time from one timezone to another
func ConvertTimezone(t time.Time, fromTZ, toTZ string) (time.Time, error) {
	fromLoc, err := time.LoadLocation(fromTZ)
	if err != nil {
		return time.Time{}, fmt.Errorf("unknown source timezone %q: %w", fromTZ, err)
	}

	toLoc, err := time.LoadLocation(toTZ)
	if err != nil {
		return time.Time{}, fmt.Errorf("unknown target timezone %q: %w", toTZ, err)
	}

	localTime := t.In(fromLoc)
	return localTime.In(toLoc), nil
}

// Common timezone abbreviations
var commonTimezones = []string{
	"UTC",
	"US/Eastern",
	"US/Central",
	"US/Mountain",
	"US/Pacific",
	"Europe/London",
	"Europe/Paris",
	"Europe/Berlin",
	"Europe/Moscow",
	"Asia/Tokyo",
	"Asia/Shanghai",
	"Asia/Kolkata",
	"Asia/Dubai",
	"Asia/Singapore",
	"Asia/Seoul",
	"Australia/Sydney",
	"Australia/Melbourne",
	"Pacific/Auckland",
	"America/Sao_Paulo",
	"America/Argentina/Buenos_Aires",
	"America/Mexico_City",
	"America/Toronto",
	"Africa/Cairo",
	"Africa/Lagos",
	"Africa/Johannesburg",
}

// ListTimezones returns common timezones
func ListTimezones() []string {
	return commonTimezones
}

// SearchTimezones finds timezones matching a pattern
func SearchTimezones(pattern string) []string {
	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return nil
	}

	var matches []string
	for _, name := range commonTimezones {
		if re.MatchString(name) {
			matches = append(matches, name)
		}
	}

	// Also check some extended timezone names
	extended := []string{
		"US/Hawaii",
		"US/Alaska",
		"Asia/Hong_Kong",
		"Asia/Bangkok",
		"Asia/Taipei",
		"Asia/Karachi",
		"Europe/Istanbul",
		"Europe/Rome",
		"Europe/Madrid",
		"Europe/Amsterdam",
		"Europe/Warsaw",
		"Europe/Zurich",
		"Europe/Stockholm",
		"Europe/Oslo",
		"Europe/Copenhagen",
		"Europe/Helsinki",
		"Europe/Athens",
		"Europe/Bucharest",
		"Europe/Kiev",
		"Africa/Nairobi",
		"Africa/Casablanca",
		"America/Chicago",
		"America/Denver",
		"America/Los_Angeles",
		"America/New_York",
		"America/Vancouver",
		"America/Halifax",
		"America/St_Johns",
		"Pacific/Honolulu",
		"Pacific/Fiji",
		"Pacific/Guam",
		"Asia/Kathmandu",
		"Asia/Kuwait",
		"Asia/Qatar",
		"Asia/Bahrain",
		"Asia/Jerusalem",
		"Asia/Beirut",
		"Asia/Tehran",
		"Asia/Kabul",
		"Asia/Tashkent",
		"Asia/Almaty",
		"Asia/Novosibirsk",
		"Asia/Vladivostok",
		"Asia/Yekaterinburg",
		"Europe/Samara",
	}
	for _, name := range extended {
		if re.MatchString(name) {
			matches = append(matches, name)
		}
	}

	// Deduplicate
	seen := make(map[string]bool)
	var unique []string
	for _, m := range matches {
		if !seen[m] {
			seen[m] = true
			unique = append(unique, m)
		}
	}
	return unique
}

// GetTimezoneInfo returns information about a timezone
func GetTimezoneInfo(tz string) (*TimezoneInfo, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("unknown timezone %q: %w", tz, err)
	}

	now := time.Now().In(loc)
	_, offset := now.Zone()

	return &TimezoneInfo{
		Name:   tz,
		Offset: offset,
		Abbrev: now.Format("MST"),
		IsDST:  now.IsDST(),
	}, nil
}

// FormatOffset formats a UTC offset as +HH:MM or -HH:MM
func FormatOffset(offset int) string {
	neg := false
	if offset < 0 {
		neg = true
		offset = -offset
	}

	hours := offset / 3600
	minutes := (offset % 3600) / 60

	sign := "+"
	if neg {
		sign = "-"
	}

	return fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)
}
