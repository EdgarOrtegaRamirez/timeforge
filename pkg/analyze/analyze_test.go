package analyze

import (
	"testing"
	"time"
)

func TestAnalyze(t *testing.T) {
	now := time.Now()
	timestamps := []time.Time{
		now.Add(-10 * time.Minute),
		now.Add(-9 * time.Minute),
		now.Add(-8 * time.Minute),
		now.Add(-1 * time.Minute), // large gap
		now,
	}

	result := Analyze(timestamps)

	if result.Count != 5 {
		t.Errorf("Count = %d, want 5", result.Count)
	}
	if result.Range != 10*time.Minute {
		t.Errorf("Range = %v, want 10m", result.Range)
	}
}

func TestAnalyzeEmpty(t *testing.T) {
	result := Analyze(nil)
	if result.Count != 0 {
		t.Errorf("Count = %d, want 0", result.Count)
	}
}

func TestAnalyzeSingle(t *testing.T) {
	now := time.Now()
	result := Analyze([]time.Time{now})
	if result.Count != 1 {
		t.Errorf("Count = %d, want 1", result.Count)
	}
}

func TestFormatReport(t *testing.T) {
	now := time.Now()
	timestamps := []time.Time{
		now.Add(-10 * time.Minute),
		now.Add(-5 * time.Minute),
		now,
	}

	result := Analyze(timestamps)
	report := FormatReport(result)

	if report == "" {
		t.Error("FormatReport() returned empty string")
	}
	if len(report) < 50 {
		t.Errorf("FormatReport() too short: %d chars", len(report))
	}
}

func TestAnalyzeWithFormats(t *testing.T) {
	entries := []TimestampEntry{
		{Timestamp: time.Now().Add(-2 * time.Minute), Raw: "2024-01-15T10:30:00Z", LineNum: 1},
		{Timestamp: time.Now().Add(-1 * time.Minute), Raw: "2024/01/15 10:31:00", LineNum: 2},
		{Timestamp: time.Now(), Raw: "2024-01-15T10:32:00Z", LineNum: 3},
	}

	result, formats := AnalyzeWithFormats(entries)

	if result.Count != 3 {
		t.Errorf("Count = %d, want 3", result.Count)
	}
	if len(formats) == 0 {
		t.Error("formats should not be empty")
	}
}
