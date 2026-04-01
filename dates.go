// This module act as a wrapper around stdlib "time"
// I am using this to avoid the complexity of timezones hours and seconds on the time.Time struct
package kalendar

import (
	"fmt"
	"time"
)

type Month int

const (
	JANUARY Month = 1 + iota
	FEBRUARY
	MARCH
	APRIL
	MAY
	JUNE
	JULY
	AUGUST
	SEPTEMBER
	OCTOBER
	NOVEMBER
	DECEMBER
)

type Weekday int

const (
	SUNDAY Weekday = iota
	MONDAY
	TUESDAY
	WEDNESDAY
	THURSDAY
	FRIDAY
	SATURDAY
)

type Date struct {
	day   int
	month Month
	year  int
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day)
}

func (d Date) Day() int     { return d.day }
func (d Date) Month() Month { return d.month }
func (d Date) Year() int    { return d.year }

func NewDate(day int, month Month, year int) Date {
	// time.Date normaliza overflow de dias, meses e anos automaticamente.
	return fromTime(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC))
}

func (d Date) toTime() time.Time {
	return time.Date(d.year, time.Month(d.month), d.day, 0, 0, 0, 0, time.UTC)
}

func fromTime(t time.Time) Date {
	return Date{day: t.Day(), month: Month(t.Month()), year: t.Year()}
}

// Plus retorna uma nova Date somando n dias.
func (d Date) Plus(days int) Date {
	return fromTime(d.toTime().AddDate(0, 0, days))
}

// Minus retorna uma nova Date subtraindo n dias.
func (d Date) Minus(days int) Date {
	return fromTime(d.toTime().AddDate(0, 0, -days))
}

// Weekday returns the day of the week (SUNDAY=0, SATURDAY=6).
func (d Date) Weekday() Weekday {
	return Weekday(d.toTime().Weekday())
}

// NextOrSame returns the nearest occurrence of the given weekday,
// starting from this date (inclusive). If d is already the target weekday,
// it returns d itself.
func (d Date) NextOrSame(w Weekday) Date {
	diff := (int(w) - int(d.Weekday()) + 7) % 7
	return d.Plus(diff)
}

// Next returns the next occurrence of the given weekday,
// strictly after this date.
func (d Date) Next(w Weekday) Date {
	return d.Plus(1).NextOrSame(w)
}

// Before returns true if d is strictly before other.
func (d Date) Before(other Date) bool {
	return d.toTime().Before(other.toTime())
}

// After returns true if d is strictly after other.
func (d Date) After(other Date) bool {
	return d.toTime().After(other.toTime())
}

type DateRange struct {
	Start Date `json:"start"`
	End   Date `json:"end"`
}

// Contains returns true if d is within the range [Start, End] inclusive.
func (r DateRange) Contains(d Date) bool {
	return !d.Before(r.Start) && !d.After(r.End)
}
