package kalendar

import "time"

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

type Date struct {
	day   int
	month Month
	year  int
}

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

