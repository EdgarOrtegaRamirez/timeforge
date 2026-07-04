package range_

import (
	"fmt"
	"time"
)

// TimeRange represents a range between two timestamps
type TimeRange struct {
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

// New creates a new TimeRange
func New(start, end time.Time) *TimeRange {
	if end.Before(start) {
		start, end = end, start
	}
	return &TimeRange{
		Start:    start,
		End:      end,
		Duration: end.Sub(start),
	}
}

// Contains checks if a time falls within the range
func (r *TimeRange) Contains(t time.Time) bool {
	return !t.Before(r.Start) && !t.After(r.End)
}

// Overlaps checks if two ranges overlap
func (r *TimeRange) Overlaps(other *TimeRange) bool {
	return r.Start.Before(other.End) && r.End.After(other.Start)
}

// Intersect returns the intersection of two ranges, or nil if no overlap
func (r *TimeRange) Intersect(other *TimeRange) *TimeRange {
	if !r.Overlaps(other) {
		return nil
	}

	start := r.Start
	if other.Start.After(start) {
		start = other.Start
	}

	end := r.End
	if other.End.Before(end) {
		end = other.End
	}

	return New(start, end)
}

// Split divides a range into equal intervals
func (r *TimeRange) Split(interval time.Duration) []*TimeRange {
	if interval <= 0 || r.Duration <= 0 {
		return nil
	}

	var ranges []*TimeRange
	current := r.Start

	for current.Before(r.End) {
		next := current.Add(interval)
		if next.After(r.End) {
			next = r.End
		}
		ranges = append(ranges, New(current, next))
		current = next
	}

	return ranges
}

// GenerateRanges creates a series of time ranges from a start time with a given interval and count
func GenerateRanges(start time.Time, interval time.Duration, count int) []*TimeRange {
	ranges := make([]*TimeRange, 0, count)
	current := start

	for i := 0; i < count; i++ {
		next := current.Add(interval)
		ranges = append(ranges, New(current, next))
		current = next
	}

	return ranges
}

// GenerateRangesUntil creates ranges from start until end time with a given interval
func GenerateRangesUntil(start time.Time, end time.Time, interval time.Duration) []*TimeRange {
	var ranges []*TimeRange
	current := start

	for current.Before(end) {
		next := current.Add(interval)
		if next.After(end) {
			next = end
		}
		ranges = append(ranges, New(current, next))
		current = next
	}

	return ranges
}

// DurationBreakdown provides a detailed breakdown of a time range
type DurationBreakdown struct {
	Years   int
	Months  int
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

// Breakdown returns a human-readable breakdown of the range duration
func (r *TimeRange) Breakdown() *DurationBreakdown {
	d := r.Duration
	if d < 0 {
		d = -d
	}

	years := int(d.Hours() / (24 * 365))
	remaining := d - time.Duration(years)*24*365*time.Hour

	months := int(remaining.Hours() / (24 * 30))
	remaining = remaining - time.Duration(months)*24*30*time.Hour

	days := int(remaining.Hours() / 24)
	remaining = remaining - time.Duration(days)*24*time.Hour

	hours := int(remaining.Hours())
	remaining = remaining - time.Duration(hours)*time.Hour

	minutes := int(remaining.Minutes())
	remaining = remaining - time.Duration(minutes)*time.Minute

	seconds := int(remaining.Seconds())

	return &DurationBreakdown{
		Years:   years,
		Months:  months,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}
}

// String returns a human-readable string of the duration breakdown
func (b *DurationBreakdown) String() string {
	var parts []string
	if b.Years > 0 {
		parts = append(parts, fmt.Sprintf("%dy", b.Years))
	}
	if b.Months > 0 {
		parts = append(parts, fmt.Sprintf("%dM", b.Months))
	}
	if b.Days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", b.Days))
	}
	if b.Hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", b.Hours))
	}
	if b.Minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", b.Minutes))
	}
	if b.Seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", b.Seconds))
	}

	result := ""
	for _, p := range parts {
		result += p
	}
	return result
}

// TimeRangeSet represents a collection of time ranges
type TimeRangeSet struct {
	Ranges []*TimeRange
}

// NewSet creates a new TimeRangeSet
func NewSet(ranges ...*TimeRange) *TimeRangeSet {
	return &TimeRangeSet{Ranges: ranges}
}

// TotalDuration returns the sum of all range durations
func (s *TimeRangeSet) TotalDuration() time.Duration {
	var total time.Duration
	for _, r := range s.Ranges {
		total += r.Duration
	}
	return total
}

// Merged returns a merged range covering the earliest start to latest end
func (s *TimeRangeSet) Merged() *TimeRange {
	if len(s.Ranges) == 0 {
		return nil
	}

	start := s.Ranges[0].Start
	end := s.Ranges[0].End

	for _, r := range s.Ranges[1:] {
		if r.Start.Before(start) {
			start = r.Start
		}
		if r.End.After(end) {
			end = r.End
		}
	}

	return New(start, end)
}
