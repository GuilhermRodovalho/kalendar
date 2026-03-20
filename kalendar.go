package kalendar

type Calendar int

const (
	GREGORIAN Calendar = iota
	JULIAN
)

type MobileDates struct {
	AshWednesday        Date `json:"ash_wednesday"`
	PalmSunday          Date `json:"palm_sunday"`
	Easter              Date `json:"easter"`
	AscensionOfTheLord  Date `json:"ascension_of_the_lord"`
	Pentecost           Date `json:"pentecost"`
	HolyTrinity         Date `json:"holy_trinity"`
	CorpusChristi       Date `json:"corpus_christi"`
	FeastOfSacredHeart  Date `json:"feast_of_the_sacred_heart"`
}

type LiturgicSeasons struct {
	Advent         DateRange `json:"advent"`
	Christmas      DateRange `json:"christmas"`
	OrdinaryTimeI  DateRange `json:"ordinary_time_i"`
	Lent           DateRange `json:"lent"`
	EasterTriduum  DateRange `json:"easter_triduum"`
	EasterSeason   DateRange `json:"easter_season"`
	OrdinaryTimeII DateRange `json:"ordinary_time_ii"`
}

type LiturgicYear struct {
	MobileDates    `json:"mobile_dates"`
	LiturgicSeasons `json:"seasons"`
}

// firstSundayOfAdvent returns the first Sunday of Advent for the given calendar year.
// It falls between November 27 and December 3.
func firstSundayOfAdvent(year int) Date {
	return NewDate(27, NOVEMBER, year).NextOrSame(SUNDAY)
}

// baptismOfTheLord returns the Feast of the Baptism of the Lord.
// It is the Sunday after Epiphany (January 6).
// If Epiphany falls on a Sunday, the Baptism is the following Monday.
func baptismOfTheLord(year int) Date {
	epiphany := NewDate(6, JANUARY, year)
	if epiphany.Weekday() == SUNDAY {
		return epiphany.Plus(1)
	}
	return epiphany.Next(SUNDAY)
}

func liturgicYearFromEaster(easter Date) *LiturgicYear {
	year := easter.year
	ashWednesday := easter.Minus(46)
	pentecost := easter.Plus(49)
	adventStart := firstSundayOfAdvent(year - 1)
	nextAdventStart := firstSundayOfAdvent(year)
	baptism := baptismOfTheLord(year)

	return &LiturgicYear{
		MobileDates: MobileDates{
			AshWednesday:       ashWednesday,
			PalmSunday:         easter.Minus(7),
			Easter:             easter,
			AscensionOfTheLord: easter.Plus(39),
			Pentecost:          pentecost,
			HolyTrinity:        easter.Plus(56),
			CorpusChristi:      easter.Plus(60),
			FeastOfSacredHeart: easter.Plus(68),
		},
		LiturgicSeasons: LiturgicSeasons{
			Advent:         DateRange{adventStart, NewDate(24, DECEMBER, year-1)},
			Christmas:      DateRange{NewDate(25, DECEMBER, year-1), baptism},
			OrdinaryTimeI:  DateRange{baptism.Plus(1), ashWednesday.Minus(1)},
			Lent:           DateRange{ashWednesday, easter.Minus(4)},
			EasterTriduum:  DateRange{easter.Minus(3), easter.Minus(1)},
			EasterSeason:   DateRange{easter, pentecost},
			OrdinaryTimeII: DateRange{pentecost.Plus(1), nextAdventStart.Minus(1)},
		},
	}
}

// This code is as implementation of Gauss algorithm to find easter day
// Available on: https://en.wikipedia.org/wiki/Date_of_Easter#Gauss's_Easter_algorithm
//
// You should pass either GREGORIAN or JULIAN as calendar, if other value is passed, GREGORIAN is considered
func EasterByGauss(year int, calendar Calendar) Date {
	a := year % 19
	b := year % 4
	c := year % 7
	m := 0
	n := 0

	switch calendar {
	case JULIAN:
		m = 15
		n = 6

	default:
		k := year / 100
		p := (13 + 8*k) / 25
		q := k / 4
		m = (15 - p + k - q) % 30
		n = (4 + k - q) % 7
	}

	d := (19*a + m) % 30
	e := (2*b + 4*c + 6*d + n) % 7
	marchEaster := d + e + 22
	aprilEaster := d + e - 9

	if aprilEaster == 25 && d == 28 && a > 10 {
		aprilEaster = 18
	}
	if aprilEaster == 26 && d == 29 && e == 6 {
		aprilEaster = 19
	}
	if marchEaster <= 31 {
		return NewDate(marchEaster, MARCH, year)
	}
	return NewDate(aprilEaster, APRIL, year)
}

// Lent returns start and end day of lent.
// Start day is Ash Wednesday, end is Holy Wednesday (day before the Easter Triduum).
func Lent(year int) (Date, Date) {
	ly := LiturgicYearOf(year)
	return ly.LiturgicSeasons.Lent.Start, ly.LiturgicSeasons.Lent.End
}

// LiturgicYearOf returns the full liturgical year for the given civil year.
func LiturgicYearOf(year int) *LiturgicYear {
	easter := EasterByGauss(year, GREGORIAN)
	return liturgicYearFromEaster(easter)
}
