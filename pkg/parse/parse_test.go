package parse

import (
	"testing"
	"time"
)

func TestParseRFC3339(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"2024-01-15T10:30:00Z", "2024-01-15T10:30:00Z"},
		{"2024-01-15T10:30:00+05:30", "2024-01-15T10:30:00+05:30"},
		{"2024-01-15T10:30:00-05:00", "2024-01-15T10:30:00-05:00"},
		{"2024-01-15T10:30:00.123456789Z", "2024-01-15T10:30:00.123456789Z"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, _, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}
			if got.Format(time.RFC3339Nano) != tt.want {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, got.Format(time.RFC3339Nano), tt.want)
			}
		})
	}
}

func TestParseDateOnly(t *testing.T) {
	tests := []struct {
		input string
		year  int
		month time.Month
		day   int
	}{
		{"2024-01-15", 2024, time.January, 15},
		{"01/15/2024", 2024, time.January, 15},
		{"15 Jan 2024", 2024, time.January, 15},
		{"Jan 15, 2024", 2024, time.January, 15},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, _, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}
			if got.Year() != tt.year || got.Month() != tt.month || got.Day() != tt.day {
				t.Errorf("Parse(%q) = %v, want %d-%02d-%02d",
					tt.input, got, tt.year, tt.month, tt.day)
			}
		})
	}
}

func TestParseUnixTimestamp(t *testing.T) {
	tests := []struct {
		input string
		unix  int64
	}{
		{"1705312200", 1705312200},
		{"1705312200000", 1705312200}, // milliseconds
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, fmt, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}
			if fmt != "unix" {
				t.Errorf("Parse(%q) format = %v, want unix", tt.input, fmt)
			}
			if got.Unix() != tt.unix {
				t.Errorf("Parse(%q) = %v, want unix %d", tt.input, got.Unix(), tt.unix)
			}
		})
	}
}

func TestParseRelative(t *testing.T) {
	before := time.Now()
	got, fmt, err := Parse("-2h30m")
	after := time.Now()

	if err != nil {
		t.Fatalf("Parse(\"-2h30m\") error = %v", err)
	}
	if fmt != "relative" {
		t.Errorf("format = %v, want relative", fmt)
	}

	expected := before.Add(-2*time.Hour - 30*time.Minute)
	if got.Before(expected.Add(-time.Second)) || got.After(after.Add(-2*time.Hour-30*time.Minute).Add(time.Second)) {
		t.Errorf("Parse(\"-2h30m\") = %v, expected around %v", got, expected)
	}
}

func TestParseEmpty(t *testing.T) {
	_, _, err := Parse("")
	if err == nil {
		t.Error("Parse(\"\") should return error")
	}
}

func TestParseInvalid(t *testing.T) {
	_, _, err := Parse("not-a-timestamp")
	if err == nil {
		t.Error("Parse(\"not-a-timestamp\") should return error")
	}
}

func TestIsTimestamp(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"2024-01-15T10:30:00Z", true},
		{"2024-01-15", true},
		{"1705312200", true},
		{"-2h", true},
		{"not-a-timestamp", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsTimestamp(tt.input); got != tt.want {
				t.Errorf("IsTimestamp(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestDetectFormat(t *testing.T) {
	got := DetectFormat("2024-01-15T10:30:00Z")
	if got == "" {
		t.Error("DetectFormat should return a non-empty format")
	}
}
