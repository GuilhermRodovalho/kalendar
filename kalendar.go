package kalendar

// CelebrationGrade representa o grau de celebração conforme o Missal Romano
type CelebrationGrade string

const (
	GradeSolenidade         CelebrationGrade = "Solenidade"
	GradeFesta              CelebrationGrade = "Festa"
	GradeMemoria            CelebrationGrade = "Memória"
	GradeMemoriaFacultativa CelebrationGrade = "Memória facultativa"
	GradeComemoracao        CelebrationGrade = "Comemoração"
)

// CelebrationLevel representa o nível de precedência (1 = maior)
type CelebrationLevel int

const (
	LevelSolenidade         CelebrationLevel = 1
	LevelFesta              CelebrationLevel = 2
	LevelMemoria            CelebrationLevel = 3
	LevelMemoriaFacultativa CelebrationLevel = 4
	LevelComemoracao        CelebrationLevel = 5
)

// CelebrationLevel retorna o nível de precedência para ordenação
func (g CelebrationGrade) Level() CelebrationLevel {
	switch g {
	case GradeSolenidade:
		return LevelSolenidade
	case GradeFesta:
		return LevelFesta
	case GradeMemoria:
		return LevelMemoria
	case GradeMemoriaFacultativa:
		return LevelMemoriaFacultativa
	case GradeComemoracao:
		return LevelComemoracao
	default:
		return LevelMemoriaFacultativa
	}
}

// Saint representa um santo do calendário
type Saint struct {
	Name             string           `json:"nome"`
	Date             string           `json:"dia"`
	Grade            CelebrationGrade `json:"grau"`
	Level            CelebrationLevel `json:"nivel"`
	Color            string           `json:"cor"`
	IsFeastOfTheLord bool             `json:"festa_do_senhor,omitempty"`
}

// Celebration representa uma celebração em uma data específica
type Celebration struct {
	Date  Date  `json:"data"`
	Saint Saint `json:"santo"`
}

// LiturgicSeasonsWithCelebrations combines liturgical seasons with celebrations
type LiturgicSeasonsWithCelebrations struct {
	MobileDates     MobileDates `json:"datas_moveis"`
	LiturgicSeasons `json:"tempos_liturgicos"`
	Celebrations    []Celebration `json:"celebrações"`
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
	Color LiturgicalColor `json:"color"`
}

type Feast struct {
	Date  Date            `json:"date"`
	Color LiturgicalColor `json:"color"`
}

type MobileDates struct {
	AshWednesday       Feast `json:"ash_wednesday"`
	PalmSunday         Feast `json:"palm_sunday"`
	GoodFriday         Feast `json:"good_friday"`
	Easter             Feast `json:"easter"`
	AscensionOfTheLord Feast `json:"ascension_of_the_lord"`
	Pentecost          Feast `json:"pentecost"`
	HolyTrinity        Feast `json:"holy_trinity"`
	CorpusChristi      Feast `json:"corpus_christi"`
	FeastOfSacredHeart Feast `json:"feast_of_the_sacred_heart"`
}

type LiturgicSeasons struct {
	Advent         Season `json:"advent"`
	Christmas      Season `json:"christmas"`
	OrdinaryTimeI  Season `json:"ordinary_time_i"`
	Lent           Season `json:"lent"`
	EasterTriduum  Season `json:"easter_triduum"`
	EasterSeason   Season `json:"easter_season"`
	OrdinaryTimeII Season `json:"ordinary_time_ii"`
}

type LiturgicYear struct {
	MobileDates     `json:"mobile_dates"`
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
			AshWednesday:       Feast{ashWednesday, Purple},
			PalmSunday:         Feast{easter.Minus(7), Red},
			GoodFriday:         Feast{easter.Minus(2), Red},
			Easter:             Feast{easter, White},
			AscensionOfTheLord: Feast{easter.Plus(39), White},
			Pentecost:          Feast{pentecost, Red},
			HolyTrinity:        Feast{easter.Plus(56), White},
			CorpusChristi:      Feast{easter.Plus(60), White},
			FeastOfSacredHeart: Feast{easter.Plus(68), White},
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
