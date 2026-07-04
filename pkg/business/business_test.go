package business

import (
	"testing"
	"time"
)

func TestIsWeekend(t *testing.T) {
	cal := NewCalendar()

	// Saturday
	sat := time.Date(2024, 1, 13, 0, 0, 0, 0, time.UTC)
	if !cal.IsWeekend(sat) {
		t.Error("Saturday should be a weekend")
	}

	// Sunday
	sun := time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC)
	if !cal.IsWeekend(sun) {
		t.Error("Sunday should be a weekend")
	}

	// Monday
	mon := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if cal.IsWeekend(mon) {
		t.Error("Monday should not be a weekend")
	}
}

func TestIsBusinessDay(t *testing.T) {
	cal := NewCalendar()

	// Monday should be a business day
	mon := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if !cal.IsBusinessDay(mon) {
		t.Error("Monday should be a business day")
	}

	// Saturday should not
	sat := time.Date(2024, 1, 13, 0, 0, 0, 0, time.UTC)
	if cal.IsBusinessDay(sat) {
		t.Error("Saturday should not be a business day")
	}
}

func TestAddBusinessDays(t *testing.T) {
	cal := NewCalendar()

	// Monday + 1 business day = Tuesday
	mon := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	tue := cal.AddBusinessDays(mon, 1)
	if tue.Weekday() != time.Tuesday {
		t.Errorf("Monday + 1 business day should be Tuesday, got %s", tue.Weekday())
	}

	// Friday + 1 business day = Monday
	fri := time.Date(2024, 1, 19, 0, 0, 0, 0, time.UTC)
	nextMon := cal.AddBusinessDays(fri, 1)
	if nextMon.Weekday() != time.Monday {
		t.Errorf("Friday + 1 business day should be Monday, got %s", nextMon.Weekday())
	}
}

func TestSubtractBusinessDays(t *testing.T) {
	cal := NewCalendar()

	// Wednesday - 2 business days = Monday
	wed := time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC)
	mon := cal.SubtractBusinessDays(wed, 2)
	if mon.Weekday() != time.Monday {
		t.Errorf("Wednesday - 2 business days should be Monday, got %s", mon.Weekday())
	}
}

func TestBusinessDaysBetween(t *testing.T) {
	cal := NewCalendar()

	// Monday to Friday = 4 business days
	mon := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	fri := time.Date(2024, 1, 19, 0, 0, 0, 0, time.UTC)
	count := cal.BusinessDaysBetween(mon, fri)
	if count != 4 {
		t.Errorf("Mon-Fri business days = %d, want 4", count)
	}

	// Monday to next Monday = 5 business days
	nextMon := time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC)
	count = cal.BusinessDaysBetween(mon, nextMon)
	if count != 5 {
		t.Errorf("Mon-nextMon business days = %d, want 5", count)
	}
}

func TestNextBusinessDay(t *testing.T) {
	cal := NewCalendar()

	// Saturday -> Monday
	sat := time.Date(2024, 1, 13, 0, 0, 0, 0, time.UTC)
	next := cal.NextBusinessDay(sat)
	if next.Weekday() != time.Monday {
		t.Errorf("Next business day after Saturday should be Monday, got %s", next.Weekday())
	}

	// Monday -> Tuesday
	mon := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	next = cal.NextBusinessDay(mon)
	if next.Weekday() != time.Tuesday {
		t.Errorf("Next business day after Monday should be Tuesday, got %s", next.Weekday())
	}
}

func TestPrevBusinessDay(t *testing.T) {
	cal := NewCalendar()

	// Monday -> Friday
	mon := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	prev := cal.PrevBusinessDay(mon)
	if prev.Weekday() != time.Friday {
		t.Errorf("Previous business day before Monday should be Friday, got %s", prev.Weekday())
	}
}

func TestBusinessDaysFrom(t *testing.T) {
	cal := NewCalendar()

	mon := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	days := cal.BusinessDaysFrom(mon, 5)

	if len(days) != 5 {
		t.Errorf("BusinessDaysFrom returned %d days, want 5", len(days))
	}

	// All should be weekdays
	for _, d := range days {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			t.Errorf("BusinessDaysFrom returned weekend day %s", d.Weekday())
		}
	}

	// Should be consecutive business days
	for i := 1; i < len(days); i++ {
		expected := cal.AddBusinessDays(days[i-1], 1)
		if !days[i].Equal(expected) {
			t.Errorf("Day %d = %v, want %v", i, days[i], expected)
		}
	}
}

func TestSummarize(t *testing.T) {
	cal := NewCalendar()

	mon := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	fri := time.Date(2024, 1, 19, 0, 0, 0, 0, time.UTC)

	summary := cal.Summarize(mon, fri)

	if summary.TotalDays != 5 {
		t.Errorf("TotalDays = %d, want 5", summary.TotalDays)
	}
	if summary.BusinessDays != 5 {
		t.Errorf("BusinessDays = %d, want 5", summary.BusinessDays)
	}
	if summary.WeekendDays != 0 {
		t.Errorf("WeekendDays = %d, want 0", summary.WeekendDays)
	}
}

func TestNewUSCalendar(t *testing.T) {
	cal := NewUSCalendar()

	// Christmas Day should be a holiday
	christmas := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
	if cal.IsBusinessDay(christmas) {
		t.Error("Christmas should not be a business day")
	}

	// July 4th should be a holiday
	july4 := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
	if cal.IsBusinessDay(july4) {
		t.Error("July 4th should not be a business day")
	}
}
