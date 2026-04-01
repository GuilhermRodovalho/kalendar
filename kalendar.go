package kalendar

// CelebrationGrade represents the celebration grade according to the Roman Missal
type CelebrationGrade string

const (
	GradeSolemnity       CelebrationGrade = "Solemnity"
	GradeFeast           CelebrationGrade = "Feast"
	GradeMemorial        CelebrationGrade = "Memorial"
	GradeOptionalMemorial CelebrationGrade = "Optional Memorial"
	GradeCommemoration   CelebrationGrade = "Commemoration"
)

// CelebrationLevel represents the precedence level (1 = highest)
type CelebrationLevel int

const (
	LevelSolemnity       CelebrationLevel = 1
	LevelFeast           CelebrationLevel = 2
	LevelMemorial        CelebrationLevel = 3
	LevelOptionalMemorial CelebrationLevel = 4
	LevelCommemoration   CelebrationLevel = 5
)

// Level returns the precedence level for sorting
func (g CelebrationGrade) Level() CelebrationLevel {
	switch g {
	case GradeSolemnity:
		return LevelSolemnity
	case GradeFeast:
		return LevelFeast
	case GradeMemorial:
		return LevelMemorial
	case GradeOptionalMemorial:
		return LevelOptionalMemorial
	case GradeCommemoration:
		return LevelCommemoration
	default:
		return LevelOptionalMemorial
	}
}

// Celebration represents a liturgical celebration on a specific date
type Celebration struct {
	Name             string           `json:"name"`
	Date             Date             `json:"date"`
	Grade            CelebrationGrade `json:"grade"`
	Level            CelebrationLevel `json:"level"`
	Color            LiturgicalColor  `json:"color"`
	IsFeastOfTheLord bool             `json:"is_feast_of_the_lord,omitempty"`
	IsMovable        bool             `json:"is_movable"`
}

// LiturgicSeasonsWithCelebrations combines liturgical seasons with celebrations
type LiturgicSeasonsWithCelebrations struct {
	LiturgicSeasons `json:"liturgical_seasons"`
	Celebrations    []Celebration `json:"celebrations"`
}

type Calendar int

const (
	GREGORIAN Calendar = iota
	JULIAN
)

type LiturgicalColor string

const (
	Green  LiturgicalColor = "green"
	Purple LiturgicalColor = "purple"
	White  LiturgicalColor = "white"
	Red    LiturgicalColor = "red"
	Rose   LiturgicalColor = "rose"
	Black  LiturgicalColor = "black"
)

type Season struct {
	DateRange
	Color LiturgicalColor
}

type Feast struct {
	Date  Date
	Color LiturgicalColor
}

type MobileDates struct {
	AshWednesday       Feast
	PalmSunday         Feast
	GoodFriday         Feast
	Easter             Feast
	AscensionOfTheLord Feast
	Pentecost          Feast
	HolyTrinity        Feast
	CorpusChristi      Feast
	FeastOfSacredHeart Feast

	Epiphany              Feast
	BaptismOfTheLord      Feast
	SaintsPeterAndPaul    Feast
	AssumptionOfMary      Feast
	AllSaints             Feast
	ChristTheKing         Feast
	HolyFamily            Feast
	MaryMotherOfTheChurch Feast
	ImmaculateHeartOfMary Feast
}

type LiturgicSeasons struct {
	Advent         Season
	Christmas      Season
	OrdinaryTimeI  Season
	Lent           Season
	EasterTriduum  Season
	EasterSeason   Season
	OrdinaryTimeII Season
}

type LiturgicYear struct {
	MobileDates
	LiturgicSeasons
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
	epiph := NewDate(6, JANUARY, year)
	if epiph.Weekday() == SUNDAY {
		return epiph.Plus(1)
	}
	return epiph.Next(SUNDAY)
}

// epiphany returns the date of Epiphany of the Lord for the given year.
// In Brazil (Próprio do Brasil), Epiphany is celebrated on the Sunday
// between January 2 and January 8.
func epiphany(year int) Date {
	return NewDate(2, JANUARY, year).NextOrSame(SUNDAY)
}

// saintsPeterAndPaul returns the date for the Solemnity of Saints Peter and Paul.
// In Brazil (Próprio do Brasil), it is celebrated on the Sunday between
// June 28 and July 4.
func saintsPeterAndPaul(year int) Date {
	return NewDate(28, JUNE, year).NextOrSame(SUNDAY)
}

// assumptionOfMary returns the date for the Solemnity of the Assumption of Mary.
// In Brazil (Próprio do Brasil), if August 15 is a Sunday it is celebrated on that day;
// otherwise it is transferred to the following Sunday.
func assumptionOfMary(year int) Date {
	return NewDate(15, AUGUST, year).NextOrSame(SUNDAY)
}

// allSaints returns the date for the Solemnity of All Saints.
// In Brazil (Próprio do Brasil), if November 1 is a Sunday it is celebrated on that day;
// otherwise it is transferred to the following Sunday.
func allSaints(year int) Date {
	return NewDate(1, NOVEMBER, year).NextOrSame(SUNDAY)
}

// christTheKing returns the date of the Solemnity of Our Lord Jesus Christ,
// King of the Universe. It is the last Sunday of Ordinary Time, i.e.,
// the Sunday before the first Sunday of Advent.
func christTheKing(year int) Date {
	advent := firstSundayOfAdvent(year)
	return advent.Minus(7)
}

// holyFamily returns the date of the Feast of the Holy Family.
// It is the Sunday within the Octave of Christmas (December 26-31).
// If there is no Sunday in that range, it is celebrated on December 30.
func holyFamily(year int) Date {
	for day := 26; day <= 31; day++ {
		d := NewDate(day, DECEMBER, year)
		if d.Weekday() == SUNDAY {
			return d
		}
	}
	return NewDate(30, DECEMBER, year)
}

// maryMotherOfTheChurch returns the date of the Memorial of the
// Blessed Virgin Mary, Mother of the Church.
// It is the Monday after Pentecost.
func maryMotherOfTheChurch(easter Date) Date {
	pentecost := easter.Plus(49)
	return pentecost.Plus(1)
}

// immaculateHeartOfMary returns the date of the Memorial of the
// Immaculate Heart of the Blessed Virgin Mary.
// It is the Saturday after the second Sunday after Pentecost
// (i.e., the Saturday after the feast of the Sacred Heart).
func immaculateHeartOfMary(easter Date) Date {
	return easter.Plus(69)
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
			AshWednesday:       Feast{ashWednesday, Purple},
			PalmSunday:         Feast{easter.Minus(7), Red},
			GoodFriday:         Feast{easter.Minus(2), Red},
			Easter:             Feast{easter, White},
			AscensionOfTheLord: Feast{easter.Plus(39), White},
			Pentecost:          Feast{pentecost, Red},
			HolyTrinity:        Feast{easter.Plus(56), White},
			CorpusChristi:      Feast{easter.Plus(60), White},
			FeastOfSacredHeart: Feast{easter.Plus(68), White},

			Epiphany:              Feast{epiphany(year), White},
			BaptismOfTheLord:      Feast{baptism, White},
			SaintsPeterAndPaul:    Feast{saintsPeterAndPaul(year), Red},
			AssumptionOfMary:      Feast{assumptionOfMary(year), White},
			AllSaints:             Feast{allSaints(year), White},
			ChristTheKing:         Feast{christTheKing(year), White},
			HolyFamily:            Feast{holyFamily(year), White},
			MaryMotherOfTheChurch: Feast{maryMotherOfTheChurch(easter), White},
			ImmaculateHeartOfMary: Feast{immaculateHeartOfMary(easter), White},
		},
		LiturgicSeasons: LiturgicSeasons{
			Advent:         Season{DateRange{adventStart, NewDate(24, DECEMBER, year-1)}, Purple},
			Christmas:      Season{DateRange{NewDate(25, DECEMBER, year-1), baptism}, White},
			OrdinaryTimeI:  Season{DateRange{baptism.Plus(1), ashWednesday.Minus(1)}, Green},
			Lent:           Season{DateRange{ashWednesday, easter.Minus(4)}, Purple},
			EasterTriduum:  Season{DateRange{easter.Minus(3), easter.Minus(1)}, White},
			EasterSeason:   Season{DateRange{easter, pentecost}, White},
			OrdinaryTimeII: Season{DateRange{pentecost.Plus(1), nextAdventStart.Minus(1)}, Green},
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
