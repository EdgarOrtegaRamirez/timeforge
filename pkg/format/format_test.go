package format

import (
	"testing"
	"time"
)

func TestFormatPreset(t *testing.T) {
	t1 := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)

	tests := []struct {
		preset string
		want   string
	}{
		{"date", "2024-01-15"},
		{"time", "10:30:45"},
		{"datetime", "2024-01-15 10:30:45"},
		{"unix", "1705314645"},
		{"year-month", "2024-01"},
	}

	for _, tt := range tests {
		t.Run(tt.preset, func(t *testing.T) {
			got, err := Format(t1, tt.preset)
			if err != nil {
				t.Fatalf("Format() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Format(%v, %q) = %q, want %q", t1, tt.preset, got, tt.want)
			}
		})
	}
}

func TestFormatRelative(t *testing.T) {
	// Test that relative format produces a non-empty string for past/future times
	past := time.Now().Add(-1 * time.Hour)
	got, err := Format(past, "relative")
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}
	if got == "" {
		t.Error("Format(past, 'relative') returned empty string")
	}
	if got == "just now" {
		t.Errorf("Format(1 hour ago, 'relative') = %q, should not be 'just now'", got)
	}

	future := time.Now().Add(1 * time.Hour)
	got, err = Format(future, "relative")
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}
	if got == "" {
		t.Error("Format(future, 'relative') returned empty string")
	}
}

func TestFormatUnix(t *testing.T) {
	t1 := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)

	got, err := Format(t1, "unix")
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}
	if got != "1705314645" {
		t.Errorf("Format(unix) = %q, want %q", got, "1705314645")
	}

	got, err = Format(t1, "unix-ms")
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}
	if got != "1705314645000" {
		t.Errorf("Format(unix-ms) = %q, want %q", got, "1705314645000")
	}
}

func TestGetPreset(t *testing.T) {
	layout, ok := GetPreset("date")
	if !ok || layout != "2006-01-02" {
		t.Errorf("GetPreset('date') = %q, %v, want %q, true", layout, ok, "2006-01-02")
	}

	_, ok = GetPreset("nonexistent")
	if ok {
		t.Error("GetPreset('nonexistent') should return false")
	}
}

func TestListFormats(t *testing.T) {
	formats := ListFormats()
	if len(formats) < 10 {
		t.Errorf("ListFormats() returned %d formats, want >= 10", len(formats))
	}

	for _, name := range []string{"date", "time", "datetime", "unix", "relative"} {
		if _, ok := formats[name]; !ok {
			t.Errorf("ListFormats() missing %q", name)
		}
	}
}
