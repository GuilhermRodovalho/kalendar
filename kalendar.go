package kalendar

import "time"

const (
	GREGORIAN = iota
	JULIAN
)

type MobileDates struct {
	ash_wednesday             Date
	palm_sunday               Date
	easter                    Date
	ascension_of_the_lord     Date
	pentecost                 Date
	holy_trinity              Date
	corpus_christi            Date
	feast_of_the_sacred_heart Date
}

type LiturgicYear struct {
	MobileDates
	// todo: add liturgic times, like Lent, Advent etc.
}

func liturgicYearFromEaster(easter Date) *LiturgicYear {
	return &LiturgicYear{
		MobileDates: MobileDates{
			ash_wednesday:             easter.Minus(46),
			palm_sunday:               easter.Minus(7),
			easter:                    easter,
			ascension_of_the_lord:     easter.Plus(39),
			pentecost:                 easter.Plus(49),
			holy_trinity:              easter.Plus(56),
			corpus_christi:            easter.Plus(60),
			feast_of_the_sacred_heart: easter.Plus(68),
		},
	}
}

// This code is as implementation of Gauss algorithm to find easter day
// Available on: https://en.wikipedia.org/wiki/Date_of_Easter#Gauss's_Easter_algorithm
//
// You should pass either GREGORIAN or JULIAN as calendar, if other value is passed, GREGORIAN is considered
func EasterByGauss(year int, calendar int) time.Time {
	a := year % 19
	b := year % 4
	c := year % 7
	M := 0
	N := 0

	switch calendar {
	case JULIAN:
		M = 15
		N = 6

	default:
		k := year / 100
		p := (13 + 8*k) / 25
		q := k / 4
		M = (15 - p + k - q) % 30
		N = (4 + k - q) % 7
	}

	d := (19*a + M) % 30
	e := (2*b + 4*c + 6*d + N) % 7
	march_easter := d + e + 22
	april_easter := d + e - 9

	if april_easter == 25 && d == 28 && a > 10 {
		april_easter = 18
	}
	if april_easter == 26 && d == 29 && e == 6 {
		april_easter = 19
	}
	if march_easter <= 31 {
		return time.Date(year, time.March, march_easter, 0, 0, 0, 0, time.UTC)
	}
	return time.Date(year, time.April, april_easter, 0, 0, 0, 0, time.UTC)
}

// returns start and end day of lent
// start day will be ash wednesday
// end will be Maundy thursday
func GetLent(year int) (time.Time, time.Time) {
	easter := EasterByGauss(year, GREGORIAN)
	ash_wednesday := easter.AddDate(0, 0, -46) // Easter - 46 days == Ash Wednesday

	return ash_wednesday, easter.AddDate(0, 0, -3) // Lent ends on Thursday
}

func GetLiturgicYearOf(year int) *LiturgicYear {
	easter := EasterByGauss(year, GREGORIAN)
	return liturgicYearFromEaster(Date{day: easter.Day(), month: Month(easter.Month()), year: easter.Year()})
}
