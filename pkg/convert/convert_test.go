package convert

import (
	"testing"
	"time"
)

func TestConvertTimezone(t *testing.T) {
	// Create a time in UTC
	utcTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	// Convert to US Eastern
	eastern, err := ConvertTimezone(utcTime, "UTC", "US/Eastern")
	if err != nil {
		t.Fatalf("ConvertTimezone() error = %v", err)
	}

	// Eastern is UTC-5 in January
	if eastern.Hour() != 5 {
		t.Errorf("Eastern hour = %d, want 5", eastern.Hour())
	}
}

func TestConvertTimezoneInvalid(t *testing.T) {
	utcTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	_, err := ConvertTimezone(utcTime, "UTC", "Invalid/Timezone")
	if err == nil {
		t.Error("ConvertTimezone() with invalid timezone should return error")
	}

	_, err = ConvertTimezone(utcTime, "Invalid/Timezone", "US/Eastern")
	if err == nil {
		t.Error("ConvertTimezone() with invalid source timezone should return error")
	}
}

func TestSearchTimezones(t *testing.T) {
	results := SearchTimezones("Eastern")
	if len(results) == 0 {
		t.Error("SearchTimezones('Eastern') returned no results")
	}

	// Check that results contain expected timezones
	found := false
	for _, tz := range results {
		if tz == "US/Eastern" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("SearchTimezones('Eastern') should contain US/Eastern, got %v", results)
	}
}

func TestGetTimezoneInfo(t *testing.T) {
	info, err := GetTimezoneInfo("US/Eastern")
	if err != nil {
		t.Fatalf("GetTimezoneInfo() error = %v", err)
	}

	if info.Name != "US/Eastern" {
		t.Errorf("Name = %q, want %q", info.Name, "US/Eastern")
	}
	if info.Offset == 0 {
		t.Error("Offset should not be 0 for US/Eastern")
	}
}

func TestFormatOffset(t *testing.T) {
	tests := []struct {
		offset int
		want   string
	}{
		{0, "+00:00"},
		{3600, "+01:00"},
		{-3600, "-01:00"},
		{19800, "+05:30"},
		{-18000, "-05:00"},
	}

	for _, tt := range tests {
		got := FormatOffset(tt.offset)
		if got != tt.want {
			t.Errorf("FormatOffset(%d) = %q, want %q", tt.offset, got, tt.want)
		}
	}
}

func TestListTimezones(t *testing.T) {
	timezones := ListTimezones()
	if len(timezones) < 10 {
		t.Errorf("ListTimezones() returned %d timezones, want >= 10", len(timezones))
	}
}
