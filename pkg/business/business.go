package business

import (
	"fmt"
	"time"
)

// Holiday represents a date that is not a business day
type Holiday struct {
	Date  time.Time
	Name  string
	Recur bool // Whether this holiday recurs annually
}

// BusinessCalendar manages business day calculations
type BusinessCalendar struct {
	weekends   map[time.Weekday]bool
	holidays   []Holiday
	holidayMap map[string]bool // "YYYY-MM-DD" -> is holiday
}

// NewCalendar creates a new business calendar with Saturday/Sunday as weekends
func NewCalendar() *BusinessCalendar {
	return &BusinessCalendar{
		weekends: map[time.Weekday]bool{
			time.Saturday: true,
			time.Sunday:   true,
		},
		holidays:   []Holiday{},
		holidayMap: make(map[string]bool),
	}
}

// NewUSCalendar creates a calendar with common US holidays
func NewUSCalendar() *BusinessCalendar {
	cal := NewCalendar()

	usHolidays := []Holiday{
		{Date: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC), Name: "New Year's Day", Recur: true},
		{Date: time.Date(0, 7, 4, 0, 0, 0, 0, time.UTC), Name: "Independence Day", Recur: true},
		{Date: time.Date(0, 12, 25, 0, 0, 0, 0, time.UTC), Name: "Christmas Day", Recur: true},
		{Date: time.Date(0, 11, 11, 0, 0, 0, 0, time.UTC), Name: "Veterans Day", Recur: true},
	}

	for _, h := range usHolidays {
		cal.AddHoliday(h)
	}

	return cal
}

// AddHoliday adds a holiday to the calendar
func (c *BusinessCalendar) AddHoliday(h Holiday) {
	c.holidays = append(c.holidays, h)
	c.holidayMap[h.Date.Format("2006-01-02")] = true
}

// IsWeekend checks if a date falls on a weekend
func (c *BusinessCalendar) IsWeekend(t time.Time) bool {
	return c.weekends[t.Weekday()]
}

// IsHoliday checks if a date is a holiday
func (c *BusinessCalendar) IsHoliday(t time.Time) bool {
	key := t.Format("2006-01-02")
	if c.holidayMap[key] {
		return true
	}

	// Check recurring holidays
	for _, h := range c.holidays {
		if h.Recur && h.Date.Month() == t.Month() && h.Date.Day() == t.Day() {
			return true
		}
	}

	return false
}

// IsBusinessDay checks if a date is a business day
func (c *BusinessCalendar) IsBusinessDay(t time.Time) bool {
	return !c.IsWeekend(t) && !c.IsHoliday(t)
}

// AddBusinessDays adds N business days to a date
func (c *BusinessCalendar) AddBusinessDays(t time.Time, days int) time.Time {
	result := t
	remaining := days

	if remaining > 0 {
		for remaining > 0 {
			result = result.AddDate(0, 0, 1)
			if c.IsBusinessDay(result) {
				remaining--
			}
		}
	} else {
		for remaining < 0 {
			result = result.AddDate(0, 0, -1)
			if c.IsBusinessDay(result) {
				remaining++
			}
		}
	}

	return result
}

// SubtractBusinessDays subtracts N business days from a date
func (c *BusinessCalendar) SubtractBusinessDays(t time.Time, days int) time.Time {
	return c.AddBusinessDays(t, -days)
}

// BusinessDaysBetween counts the number of business days between two dates
func (c *BusinessCalendar) BusinessDaysBetween(start, end time.Time) int {
	if start.After(end) {
		start, end = end, start
	}

	count := 0
	current := start

	for current.Before(end) {
		current = current.AddDate(0, 0, 1)
		if c.IsBusinessDay(current) {
			count++
		}
	}

	return count
}

// NextBusinessDay returns the next business day from a given date
func (c *BusinessCalendar) NextBusinessDay(t time.Time) time.Time {
	next := t.AddDate(0, 0, 1)
	for !c.IsBusinessDay(next) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

// PrevBusinessDay returns the previous business day from a given date
func (c *BusinessCalendar) PrevBusinessDay(t time.Time) time.Time {
	prev := t.AddDate(0, 0, -1)
	for !c.IsBusinessDay(prev) {
		prev = prev.AddDate(0, 0, -1)
	}
	return prev
}

// BusinessDaysFrom returns a slice of business days starting from a date
func (c *BusinessCalendar) BusinessDaysFrom(start time.Time, count int) []time.Time {
	days := make([]time.Time, 0, count)
	current := start

	for len(days) < count {
		if c.IsBusinessDay(current) {
			days = append(days, current)
		}
		current = current.AddDate(0, 0, 1)
	}

	return days
}

// FormatBusinessDays formats a list of business days
func FormatBusinessDays(days []time.Time, layout string) []string {
	result := make([]string, len(days))
	for i, d := range days {
		result[i] = d.Format(layout)
	}
	return result
}

// BusinessDaySummary provides a summary of business day calculations
type BusinessDaySummary struct {
	StartDate      time.Time
	EndDate        time.Time
	TotalDays      int
	BusinessDays   int
	WeekendDays    int
	HolidayDays    int
	FormattedRange string
}

// Summarize returns a summary of business days in a range
func (c *BusinessCalendar) Summarize(start, end time.Time) *BusinessDaySummary {
	if start.After(end) {
		start, end = end, start
	}

	total := 0
	business := 0
	weekends := 0
	holidays := 0

	current := start
	for !current.After(end) {
		total++
		if c.IsWeekend(current) {
			weekends++
		} else if c.IsHoliday(current) {
			holidays++
		} else {
			business++
		}
		current = current.AddDate(0, 0, 1)
	}

	return &BusinessDaySummary{
		StartDate:      start,
		EndDate:        end,
		TotalDays:      total,
		BusinessDays:   business,
		WeekendDays:    weekends,
		HolidayDays:    holidays,
		FormattedRange: fmt.Sprintf("%s to %s", start.Format("2006-01-02"), end.Format("2006-01-02")),
	}
}
