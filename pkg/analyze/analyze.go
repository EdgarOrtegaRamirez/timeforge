package analyze

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// TimestampEntry represents a parsed timestamp from a log
type TimestampEntry struct {
	Timestamp time.Time
	Raw       string
	LineNum   int
}

// AnalysisResult contains the analysis of timestamp data
type AnalysisResult struct {
	Count       int
	Min         time.Time
	Max         time.Time
	Range       time.Duration
	Mean        time.Time
	Median      time.Time
	P50         time.Duration
	P95         time.Duration
	P99         time.Duration
	StdDev      time.Duration
	Gaps        []TimeGap
	Buckets     []TimeBucket
	Bursts      []TimeBurst
	FormatStats map[string]int
}

// TimeGap represents a gap between consecutive timestamps
type TimeGap struct {
	Start time.Time
	End   time.Time
	Gap   time.Duration
}

// TimeBucket represents a count of timestamps in a time interval
type TimeBucket struct {
	Start time.Time
	End   time.Time
	Count int
}

// TimeBurst represents a period of high activity
type TimeBurst struct {
	Start    time.Time
	End      time.Time
	Count    int
	AvgGap   time.Duration
	Severity string // "low", "medium", "high"
}

// Analyze performs comprehensive analysis on a list of timestamps
func Analyze(timestamps []time.Time) *AnalysisResult {
	if len(timestamps) == 0 {
		return &AnalysisResult{}
	}

	// Sort timestamps
	sorted := make([]time.Time, len(timestamps))
	copy(sorted, timestamps)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Before(sorted[j])
	})

	result := &AnalysisResult{
		Count: len(sorted),
		Min:   sorted[0],
		Max:   sorted[len(sorted)-1],
		Range: sorted[len(sorted)-1].Sub(sorted[0]),
	}

	// Calculate mean
	var total time.Duration
	for _, t := range sorted {
		total += t.Sub(sorted[0])
	}
	result.Mean = sorted[0].Add(total / time.Duration(len(sorted)))

	// Calculate median
	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		result.Median = sorted[mid-1].Add(sorted[mid].Sub(sorted[mid-1]) / 2)
	} else {
		result.Median = sorted[mid]
	}

	// Calculate percentiles
	gaps := make([]time.Duration, 0, len(sorted)-1)
	for i := 1; i < len(sorted); i++ {
		gaps = append(gaps, sorted[i].Sub(sorted[i-1]))
	}
	sort.Slice(gaps, func(i, j int) bool { return gaps[i] < gaps[j] })

	if len(gaps) > 0 {
		result.P50 = percentile(gaps, 50)
		result.P95 = percentile(gaps, 95)
		result.P99 = percentile(gaps, 99)
	}

	// Calculate standard deviation
	if len(gaps) > 0 {
		var sum float64
		mean := float64(result.Range) / float64(len(sorted)-1)
		for _, g := range gaps {
			diff := float64(g) - mean
			sum += diff * diff
		}
		result.StdDev = time.Duration(math.Sqrt(sum / float64(len(gaps))))
	}

	// Find gaps
	result.Gaps = findGaps(sorted, result.P95*3)

	// Create buckets
	result.Buckets = createBuckets(sorted, result.Range)

	// Find bursts
	result.Bursts = findBursts(sorted, result.P50/2)

	return result
}

// AnalyzeWithFormats also tracks which timestamp formats were used
func AnalyzeWithFormats(entries []TimestampEntry) (*AnalysisResult, map[string]int) {
	timestamps := make([]time.Time, len(entries))
	for i, e := range entries {
		timestamps[i] = e.Timestamp
	}

	result := Analyze(timestamps)

	// Count formats
	formatCounts := make(map[string]int)
	for _, e := range entries {
		detected := detectFormat(e.Raw)
		formatCounts[detected]++
	}

	return result, formatCounts
}

// FormatReport generates a text report of the analysis
func FormatReport(r *AnalysisResult) string {
	var sb strings.Builder

	sb.WriteString("Timestamp Analysis Report\n")
	sb.WriteString(strings.Repeat("=", 40) + "\n\n")

	sb.WriteString(fmt.Sprintf("Count:     %d\n", r.Count))
	sb.WriteString(fmt.Sprintf("Time Span: %s\n", formatDuration(r.Range)))
	sb.WriteString(fmt.Sprintf("Earliest:  %s\n", r.Min.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Latest:    %s\n", r.Max.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Mean:      %s\n", r.Mean.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Median:    %s\n", r.Median.Format(time.RFC3339)))

	if r.Count > 1 {
		sb.WriteString("\nGap Statistics\n")
		sb.WriteString(strings.Repeat("-", 40) + "\n")
		sb.WriteString(fmt.Sprintf("P50 Gap:   %s\n", formatDuration(r.P50)))
		sb.WriteString(fmt.Sprintf("P95 Gap:   %s\n", formatDuration(r.P95)))
		sb.WriteString(fmt.Sprintf("P99 Gap:   %s\n", formatDuration(r.P99)))
		sb.WriteString(fmt.Sprintf("Std Dev:   %s\n", formatDuration(r.StdDev)))
	}

	if len(r.Gaps) > 0 {
		sb.WriteString(fmt.Sprintf("\nSignificant Gaps: %d\n", len(r.Gaps)))
		sb.WriteString(strings.Repeat("-", 40) + "\n")
		for i, g := range r.Gaps {
			if i >= 5 {
				sb.WriteString(fmt.Sprintf("  ... and %d more\n", len(r.Gaps)-5))
				break
			}
			sb.WriteString(fmt.Sprintf("  %s → %s (%s)\n",
				g.Start.Format("15:04:05"),
				g.End.Format("15:04:05"),
				formatDuration(g.Gap)))
		}
	}

	if len(r.Bursts) > 0 {
		sb.WriteString(fmt.Sprintf("\nActivity Bursts: %d\n", len(r.Bursts)))
		sb.WriteString(strings.Repeat("-", 40) + "\n")
		for i, b := range r.Bursts {
			if i >= 5 {
				sb.WriteString(fmt.Sprintf("  ... and %d more\n", len(r.Bursts)-5))
				break
			}
			sb.WriteString(fmt.Sprintf("  [%s] %d events, avg gap %s (%s)\n",
				b.Severity,
				b.Count,
				formatDuration(b.AvgGap),
				b.Start.Format("15:04:05")))
		}
	}

	return sb.String()
}

func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(math.Ceil(p/100*float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

func findGaps(timestamps []time.Time, threshold time.Duration) []TimeGap {
	var gaps []TimeGap
	for i := 1; i < len(timestamps); i++ {
		gap := timestamps[i].Sub(timestamps[i-1])
		if gap > threshold {
			gaps = append(gaps, TimeGap{
				Start: timestamps[i-1],
				End:   timestamps[i],
				Gap:   gap,
			})
		}
	}
	return gaps
}

func createBuckets(timestamps []time.Time, totalRange time.Duration) []TimeBucket {
	if len(timestamps) == 0 || totalRange == 0 {
		return nil
	}

	// Create ~20 buckets
	bucketCount := 20
	if bucketCount > len(timestamps) {
		bucketCount = len(timestamps)
	}

	bucketSize := totalRange / time.Duration(bucketCount)
	buckets := make([]TimeBucket, bucketCount)

	for i := 0; i < bucketCount; i++ {
		start := timestamps[0].Add(time.Duration(i) * bucketSize)
		end := start.Add(bucketSize)
		buckets[i] = TimeBucket{Start: start, End: end}
	}

	// Count timestamps in each bucket
	for _, ts := range timestamps {
		idx := int(ts.Sub(timestamps[0]) / bucketSize)
		if idx >= bucketCount {
			idx = bucketCount - 1
		}
		buckets[idx].Count++
	}

	return buckets
}

func findBursts(timestamps []time.Time, avgGapThreshold time.Duration) []TimeBurst {
	if len(timestamps) < 3 {
		return nil
	}

	var bursts []TimeBurst
	current := TimeBurst{
		Start: timestamps[0],
	}

	count := 0
	var totalGap time.Duration

	for i := 1; i < len(timestamps); i++ {
		gap := timestamps[i].Sub(timestamps[i-1])

		if gap < avgGapThreshold {
			count++
			totalGap += gap
			current.End = timestamps[i]
			current.Count = count
		} else {
			if count >= 3 {
				current.AvgGap = totalGap / time.Duration(count)
				current.Severity = classifySeverity(current.AvgGap, avgGapThreshold)
				bursts = append(bursts, current)
			}
			current = TimeBurst{
				Start: timestamps[i],
			}
			count = 0
			totalGap = 0
		}
	}

	// Check final burst
	if count >= 3 {
		current.AvgGap = totalGap / time.Duration(count)
		current.Severity = classifySeverity(current.AvgGap, avgGapThreshold)
		bursts = append(bursts, current)
	}

	return bursts
}

func classifySeverity(avgGap, threshold time.Duration) string {
	ratio := float64(threshold) / float64(avgGap)
	switch {
	case ratio > 10:
		return "high"
	case ratio > 3:
		return "medium"
	default:
		return "low"
	}
}

func detectFormat(raw string) string {
	if strings.Contains(raw, "T") && strings.Contains(raw, ":") {
		return "ISO8601"
	}
	if strings.Contains(raw, "/") && strings.Contains(raw, ":") {
		return "log-format"
	}
	if strings.Count(raw, "-") >= 2 {
		return "date-dash"
	}
	return "unknown"
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.2fµs", float64(d.Microseconds()))
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Milliseconds()))
	}
	if d < time.Minute {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}
