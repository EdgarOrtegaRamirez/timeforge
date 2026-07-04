package duration

import (
	"testing"
	"time"
)

func TestFromDuration(t *testing.T) {
	d := 2*time.Hour + 30*time.Minute + 15*time.Second
	result := FromDuration(d)

	if result.Seconds != 9015.0 {
		t.Errorf("Seconds = %v, want 9015", result.Seconds)
	}
	if result.Minutes != 150.25 {
		t.Errorf("Minutes = %v, want 150.25", result.Minutes)
	}
	if result.Hours != 2.5041666666666665 {
		t.Errorf("Hours = %v, want ~2.5", result.Hours)
	}
	if result.Breakdown.Days != 0 || result.Breakdown.Hours != 2 || result.Breakdown.Minutes != 30 || result.Breakdown.Seconds != 15 {
		t.Errorf("Breakdown = %+v, want {0 2 30 15 0}", result.Breakdown)
	}
}

func TestFormatCompact(t *testing.T) {
	tests := []struct {
		input time.Duration
		want  string
	}{
		{0, "0s"},
		{5 * time.Second, "5s"},
		{2*time.Minute + 30*time.Second, "2m30s"},
		{1*time.Hour + 30*time.Minute, "1h30m"},
		{24*time.Hour + 5*time.Hour, "1d5h"},
		{-2*time.Hour - 30*time.Minute, "-2h30m"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := FormatCompact(tt.input)
			if got != tt.want {
				t.Errorf("FormatCompact(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatHuman(t *testing.T) {
	tests := []struct {
		input time.Duration
		want  string
	}{
		{0, "zero"},
		{5 * time.Second, "5 seconds"},
		{time.Second, "1 second"},
		{2*time.Hour + 30*time.Minute, "2 hours, 30 minutes"},
		{-1 * time.Hour, "-1 hour"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := FormatHuman(tt.input)
			if got != tt.want {
				t.Errorf("FormatHuman(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input string
		want  time.Duration
	}{
		{"2h30m", 2*time.Hour + 30*time.Minute},
		{"1d", 24 * time.Hour},
		{"1w", 7 * 24 * time.Hour},
		{"1M", 30 * 24 * time.Hour},
		{"1y", 365 * 24 * time.Hour},
		{"500ms", 500 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	d1 := 1 * time.Hour
	d2 := 30 * time.Minute
	d3 := 15 * time.Second

	got := Add(d1, d2, d3)
	want := 1*time.Hour + 30*time.Minute + 15*time.Second

	if got != want {
		t.Errorf("Add() = %v, want %v", got, want)
	}
}

func TestSubtract(t *testing.T) {
	base := 2 * time.Hour
	sub := []time.Duration{30 * time.Minute, 15 * time.Second}

	got := Subtract(base, sub...)
	want := 1*time.Hour + 29*time.Minute + 45*time.Second

	if got != want {
		t.Errorf("Subtract() = %v, want %v", got, want)
	}
}

func TestMin(t *testing.T) {
	got := Min(3*time.Hour, 1*time.Hour, 2*time.Hour)
	want := 1 * time.Hour

	if got != want {
		t.Errorf("Min() = %v, want %v", got, want)
	}
}

func TestMax(t *testing.T) {
	got := Max(3*time.Hour, 1*time.Hour, 2*time.Hour)
	want := 3 * time.Hour

	if got != want {
		t.Errorf("Max() = %v, want %v", got, want)
	}
}

func TestAverage(t *testing.T) {
	got := Average(1*time.Hour, 2*time.Hour, 3*time.Hour)
	want := 2 * time.Hour

	if got != want {
		t.Errorf("Average() = %v, want %v", got, want)
	}
}

func TestDifference(t *testing.T) {
	got := Difference(3*time.Hour, 1*time.Hour)
	want := 2 * time.Hour

	if got != want {
		t.Errorf("Difference() = %v, want %v", got, want)
	}

	got = Difference(1*time.Hour, 3*time.Hour)
	if got != want {
		t.Errorf("Difference() = %v, want %v", got, want)
	}
}

func TestPercentage(t *testing.T) {
	got := Percentage(100*time.Second, 25*time.Second)
	if got != 25.0 {
		t.Errorf("Percentage() = %v, want 25.0", got)
	}
}

func TestParseEmpty(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Error("Parse(\"\") should return error")
	}
}
