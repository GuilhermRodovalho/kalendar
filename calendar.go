package kalendar

import "sort"

// SeasonName identifies a liturgical season.
type SeasonName string

const (
	SeasonAdvent         SeasonName = "advent"
	SeasonChristmas      SeasonName = "christmas"
	SeasonOrdinaryTimeI  SeasonName = "ordinary_time_i"
	SeasonLent           SeasonName = "lent"
	SeasonEasterTriduum  SeasonName = "easter_triduum"
	SeasonEasterSeason   SeasonName = "easter_season"
	SeasonOrdinaryTimeII SeasonName = "ordinary_time_ii"
)

// DayCelebration represents a celebration without its date,
// since the date is already present in the parent CalendarEntry.
type DayCelebration struct {
	Name             string           `json:"name"`
	Grade            CelebrationGrade `json:"grade"`
	Level            CelebrationLevel `json:"level"`
	Color            LiturgicalColor  `json:"color"`
	IsFeastOfTheLord bool             `json:"is_feast_of_the_lord,omitempty"`
	IsMovable        bool             `json:"is_movable"`
	MissalPage       int              `json:"missal_page,omitempty"`
	Prefaces         []Preface        `json:"prefaces,omitempty"`
	CommonRef        *CommonReference `json:"common_ref,omitempty"`
}

// CalendarEntry represents a single day in the liturgical calendar.
type CalendarEntry struct {
	Date             Date             `json:"date"`
	Season           SeasonName       `json:"season"`
	SeasonColor      LiturgicalColor  `json:"season_color"`
	Celebrations     []DayCelebration `json:"celebrations"`
	SeasonMissalPage int              `json:"season_missal_page,omitempty"`
	SeasonPrefaces   []Preface        `json:"season_prefaces,omitempty"`
}

type namedSeason struct {
	name   SeasonName
	season Season
}

// seasonOrder returns the ordered list of seasons for a liturgical year.
// Allocated once per calendar operation instead of per-day.
func seasonOrder(seasons LiturgicSeasons) [7]namedSeason {
	return [7]namedSeason{
		{SeasonAdvent, seasons.Advent},
		{SeasonChristmas, seasons.Christmas},
		{SeasonOrdinaryTimeI, seasons.OrdinaryTimeI},
		{SeasonLent, seasons.Lent},
		{SeasonEasterTriduum, seasons.EasterTriduum},
		{SeasonEasterSeason, seasons.EasterSeason},
		{SeasonOrdinaryTimeII, seasons.OrdinaryTimeII},
	}
}

// seasonForDate returns the season name, color, and which season set index matched.
// It checks multiple season sets (current and next liturgical year) so that
// the Advent/Christmas seasons at the end of the civil year are found.
func seasonForDate(d Date, seasonSets ...*[7]namedSeason) (SeasonName, LiturgicalColor, int) {
	for idx, ordered := range seasonSets {
		for _, ns := range ordered {
			if ns.season.DateRange.Contains(d) {
				return ns.name, ns.season.Color, idx
			}
		}
	}

	return SeasonOrdinaryTimeII, Green, 0
}

// toDayCelebration converts a Celebration to a DayCelebration (drops the Date).
func toDayCelebration(c Celebration) DayCelebration {
	return DayCelebration{
		Name:             c.Name,
		Grade:            c.Grade,
		Level:            c.Level,
		Color:            c.Color,
		IsFeastOfTheLord: c.IsFeastOfTheLord,
		IsMovable:        c.IsMovable,
	}
}

// buildCelebrationIndex groups celebrations by their date string for O(1) lookup.
func buildCelebrationIndex(celebrations []Celebration) map[Date][]DayCelebration {
	index := make(map[Date][]DayCelebration, len(celebrations))
	for _, c := range celebrations {
		index[c.Date] = append(index[c.Date], toDayCelebration(c))
	}
	// Sort each day's celebrations by precedence (lower Level = higher precedence)
	for date := range index {
		sort.Slice(index[date], func(i, j int) bool {
			return index[date][i].Level < index[date][j].Level
		})
	}
	return index
}

// daysInYear returns the number of days in a given year.
func daysInYear(year int) int {
	start := NewDate(1, JANUARY, year)
	end := NewDate(1, JANUARY, year+1)
	return int(end.toTime().Sub(start.toTime()).Hours() / 24)
}

// GetCalendar returns a CalendarEntry for each day of the given civil year.
func GetCalendar(year int) ([]CalendarEntry, error) {
	ly := LiturgicYearOf(year)
	nextLy := LiturgicYearOf(year + 1)
	celebrations, err := GetCelebrationsForYear(year)
	if err != nil {
		return nil, err
	}

	index := buildCelebrationIndex(celebrations)
	days := daysInYear(year)
	entries := make([]CalendarEntry, days)
	ordered := seasonOrder(ly.LiturgicSeasons)
	nextOrdered := seasonOrder(nextLy.LiturgicSeasons)
	startDate := NewDate(1, JANUARY, year)

	lyByIndex := []*LiturgicYear{ly, nextLy}

	for i := range days {
		d := startDate.Plus(i)
		seasonName, seasonColor, lyIdx := seasonForDate(d, &ordered, &nextOrdered)

		dayCelebrations := index[d]
		if dayCelebrations == nil {
			dayCelebrations = []DayCelebration{}
		}

		seasonPage, seasonPrefaces := resolveSeasonMissal(d, seasonName, lyByIndex[lyIdx])

		for i := range dayCelebrations {
			ref := resolveCelebrationMissal(dayCelebrations[i].Name, d, dayCelebrations[i].IsMovable)
			if ref != nil {
				dayCelebrations[i].MissalPage = ref.MissalPage
				dayCelebrations[i].Prefaces = ref.Prefaces
				dayCelebrations[i].CommonRef = ref.CommonRef
			}
		}

		entries[i] = CalendarEntry{
			Date:             d,
			Season:           seasonName,
			SeasonColor:      seasonColor,
			Celebrations:     dayCelebrations,
			SeasonMissalPage: seasonPage,
			SeasonPrefaces:   seasonPrefaces,
		}
	}

	return entries, nil
}

// GetCalendarDay returns the CalendarEntry for a specific date.
func GetCalendarDay(year int, month Month, day int) (*CalendarEntry, error) {
	ly := LiturgicYearOf(year)
	nextLy := LiturgicYearOf(year + 1)
	celebrations, err := GetCelebrationsForYear(year)
	if err != nil {
		return nil, err
	}

	d := NewDate(day, month, year)
	ordered := seasonOrder(ly.LiturgicSeasons)
	nextOrdered := seasonOrder(nextLy.LiturgicSeasons)
	seasonName, seasonColor, lyIdx := seasonForDate(d, &ordered, &nextOrdered)

	lyByIndex := []*LiturgicYear{ly, nextLy}
	dayCelebrations := filterCelebrationsForDate(celebrations, d)

	seasonPage, seasonPrefaces := resolveSeasonMissal(d, seasonName, lyByIndex[lyIdx])

	for i := range dayCelebrations {
		ref := resolveCelebrationMissal(dayCelebrations[i].Name, d, dayCelebrations[i].IsMovable)
		if ref != nil {
			dayCelebrations[i].MissalPage = ref.MissalPage
			dayCelebrations[i].Prefaces = ref.Prefaces
			dayCelebrations[i].CommonRef = ref.CommonRef
		}
	}

	return &CalendarEntry{
		Date:             d,
		Season:           seasonName,
		SeasonColor:      seasonColor,
		Celebrations:     dayCelebrations,
		SeasonMissalPage: seasonPage,
		SeasonPrefaces:   seasonPrefaces,
	}, nil
}

// filterCelebrationsForDate returns sorted DayCelebrations matching the given date.
func filterCelebrationsForDate(celebrations []Celebration, d Date) []DayCelebration {
	var result []DayCelebration
	for _, c := range celebrations {
		if c.Date == d {
			result = append(result, toDayCelebration(c))
		}
	}
	if result == nil {
		return []DayCelebration{}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Level < result[j].Level
	})
	return result
}

// GetMobileDates returns only the mobile celebrations for a given year.
func GetMobileDates(year int) []Celebration {
	ly := LiturgicYearOf(year)
	return loadMobileCelebrations(ly)
}
