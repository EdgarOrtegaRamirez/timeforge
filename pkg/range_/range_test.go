package range_

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tr := New(start, end)

	if !tr.Start.Equal(start) {
		t.Errorf("Start = %v, want %v", tr.Start, start)
	}
	if !tr.End.Equal(end) {
		t.Errorf("End = %v, want %v", tr.End, end)
	}
	if tr.Duration != 2*time.Hour {
		t.Errorf("Duration = %v, want %v", tr.Duration, 2*time.Hour)
	}
}

func TestNewReversed(t *testing.T) {
	start := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	tr := New(start, end)

	if !tr.Start.Equal(end) {
		t.Errorf("Start = %v, want %v", tr.Start, end)
	}
	if !tr.End.Equal(start) {
		t.Errorf("End = %v, want %v", tr.End, start)
	}
}

func TestContains(t *testing.T) {
	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	tr := New(start, end)

	tests := []struct {
		time time.Time
		want bool
	}{
		{time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC), true},
		{time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), true},
		{time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC), true},
		{time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC), false},
		{time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC), false},
	}

	for _, tt := range tests {
		if got := tr.Contains(tt.time); got != tt.want {
			t.Errorf("Contains(%v) = %v, want %v", tt.time, got, tt.want)
		}
	}
}

func TestOverlaps(t *testing.T) {
	r1 := New(
		time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	)
	r2 := New(
		time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC),
	)
	r3 := New(
		time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC),
	)

	if !r1.Overlaps(r2) {
		t.Error("r1 and r2 should overlap")
	}
	if r1.Overlaps(r3) {
		t.Error("r1 and r3 should not overlap")
	}
}

func TestSplit(t *testing.T) {
	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	tr := New(start, end)

	ranges := tr.Split(30 * time.Minute)

	if len(ranges) != 4 {
		t.Errorf("Split() returned %d ranges, want 4", len(ranges))
	}

	for i, r := range ranges {
		if r.Duration != 30*time.Minute {
			t.Errorf("Range %d duration = %v, want 30m", i, r.Duration)
		}
	}
}

func TestGenerateRanges(t *testing.T) {
	start := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	interval := time.Hour

	ranges := GenerateRanges(start, interval, 5)

	if len(ranges) != 5 {
		t.Errorf("GenerateRanges() returned %d ranges, want 5", len(ranges))
	}

	for i, r := range ranges {
		expected := start.Add(time.Duration(i) * interval)
		if !r.Start.Equal(expected) {
			t.Errorf("Range %d start = %v, want %v", i, r.Start, expected)
		}
	}
}

func TestBreakdown(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 3, 15, 10, 30, 45, 0, time.UTC)
	tr := New(start, end)

	bd := tr.Breakdown()

	if bd.Months != 2 {
		t.Errorf("Months = %d, want 2", bd.Months)
	}
	if bd.Days != 14 {
		t.Errorf("Days = %d, want 14", bd.Days)
	}
	if bd.Hours != 10 {
		t.Errorf("Hours = %d, want 10", bd.Hours)
	}
}

func TestTimeRangeSet(t *testing.T) {
	r1 := New(
		time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	)
	r2 := New(
		time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 16, 0, 0, 0, time.UTC),
	)

	set := NewSet(r1, r2)

	if set.TotalDuration() != 4*time.Hour {
		t.Errorf("TotalDuration = %v, want 4h", set.TotalDuration())
	}

	merged := set.Merged()
	if !merged.Start.Equal(r1.Start) || !merged.End.Equal(r2.End) {
		t.Errorf("Merged range = %v, want %v to %v", merged, r1.Start, r2.End)
	}
}
