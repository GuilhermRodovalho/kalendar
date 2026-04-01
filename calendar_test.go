package kalendar

import "testing"

func TestGetCalendarLength(t *testing.T) {
	entries, err := GetCalendar(2026)
	if err != nil {
		t.Fatalf("GetCalendar(2026) error = %v", err)
	}
	if len(entries) != 365 {
		t.Errorf("expected 365 entries, got %d", len(entries))
	}
}

func TestGetCalendarLeapYear(t *testing.T) {
	entries, err := GetCalendar(2024)
	if err != nil {
		t.Fatalf("GetCalendar(2024) error = %v", err)
	}
	if len(entries) != 366 {
		t.Errorf("expected 366 entries for leap year, got %d", len(entries))
	}
}

func TestGetCalendarFirstAndLastDay(t *testing.T) {
	entries, err := GetCalendar(2026)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	first := entries[0]
	if first.Date != NewDate(1, JANUARY, 2026) {
		t.Errorf("first day = %v, want 2026-01-01", first.Date)
	}

	last := entries[len(entries)-1]
	if last.Date != NewDate(31, DECEMBER, 2026) {
		t.Errorf("last day = %v, want 2026-12-31", last.Date)
	}
}

func TestGetCalendarEveryDayHasSeason(t *testing.T) {
	entries, err := GetCalendar(2026)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	for _, e := range entries {
		if e.Season == "" {
			t.Errorf("date %v has empty season", e.Date)
		}
		if e.SeasonColor == "" {
			t.Errorf("date %v has empty season color", e.Date)
		}
	}
}

func TestGetCalendarCelebrationsNotNil(t *testing.T) {
	entries, err := GetCalendar(2026)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	for _, e := range entries {
		if e.Celebrations == nil {
			t.Errorf("date %v has nil celebrations (should be empty slice)", e.Date)
		}
	}
}

func TestGetCalendarSeasonTransitions(t *testing.T) {
	entries, err := GetCalendar(2026)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	dateMap := make(map[Date]CalendarEntry, len(entries))
	for _, e := range entries {
		dateMap[e.Date] = e
	}

	// 2026: Ash Wednesday = Feb 18 (Lent starts)
	ashWed := dateMap[NewDate(18, FEBRUARY, 2026)]
	if ashWed.Season != SeasonLent {
		t.Errorf("Ash Wednesday season = %v, want lent", ashWed.Season)
	}

	dayBeforeAshWed := dateMap[NewDate(17, FEBRUARY, 2026)]
	if dayBeforeAshWed.Season != SeasonOrdinaryTimeI {
		t.Errorf("day before Ash Wednesday season = %v, want ordinary_time_i", dayBeforeAshWed.Season)
	}

	// 2026: Easter = April 5
	easter := dateMap[NewDate(5, APRIL, 2026)]
	if easter.Season != SeasonEasterSeason {
		t.Errorf("Easter season = %v, want easter_season", easter.Season)
	}

	// 2026: Pentecost = May 24
	pentecost := dateMap[NewDate(24, MAY, 2026)]
	if pentecost.Season != SeasonEasterSeason {
		t.Errorf("Pentecost season = %v, want easter_season", pentecost.Season)
	}

	dayAfterPentecost := dateMap[NewDate(25, MAY, 2026)]
	if dayAfterPentecost.Season != SeasonOrdinaryTimeII {
		t.Errorf("day after Pentecost season = %v, want ordinary_time_ii", dayAfterPentecost.Season)
	}
}

func TestGetCalendarCelebrationsSortedByPrecedence(t *testing.T) {
	entries, err := GetCalendar(2026)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	for _, e := range entries {
		for i := 1; i < len(e.Celebrations); i++ {
			if e.Celebrations[i].Level < e.Celebrations[i-1].Level {
				t.Errorf("date %v: celebrations not sorted by precedence (%v after %v)",
					e.Date, e.Celebrations[i].Name, e.Celebrations[i-1].Name)
			}
		}
	}
}

func TestGetCalendarDaySpecific(t *testing.T) {
	// Jan 1 = Santa Maria, Mãe de Deus (Solemnity)
	entry, err := GetCalendarDay(2026, JANUARY, 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	if entry.Date != NewDate(1, JANUARY, 2026) {
		t.Errorf("date = %v, want 2026-01-01", entry.Date)
	}

	if len(entry.Celebrations) == 0 {
		t.Fatal("expected celebrations on Jan 1")
	}

	found := false
	for _, c := range entry.Celebrations {
		if c.Name == "Santa Maria, Mãe de Deus" {
			found = true
			if c.Grade != GradeSolemnity {
				t.Errorf("grade = %v, want Solemnity", c.Grade)
			}
		}
	}
	if !found {
		t.Error("Santa Maria, Mãe de Deus not found on Jan 1")
	}
}

func TestGetCalendarDayWithMobileCelebration(t *testing.T) {
	// 2026: Epiphany is Jan 4 (first Sunday between Jan 2-8)
	entry, err := GetCalendarDay(2026, JANUARY, 4)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	found := false
	for _, c := range entry.Celebrations {
		if c.Name == "Epifania do Senhor" {
			found = true
			if !c.IsMovable {
				t.Error("Epiphany should be movable")
			}
		}
	}
	if !found {
		t.Error("Epifania do Senhor not found on Jan 4 2026")
	}
}

func TestGetCalendarDayEmptyCelebrations(t *testing.T) {
	// Pick a day unlikely to have celebrations
	entry, err := GetCalendarDay(2026, MARCH, 10)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	if entry.Celebrations == nil {
		t.Error("celebrations should be empty slice, not nil")
	}
}

func TestGetMobileDates(t *testing.T) {
	mobile := GetMobileDates(2026)

	if len(mobile) == 0 {
		t.Fatal("GetMobileDates should return celebrations")
	}

	for _, c := range mobile {
		if !c.IsMovable {
			t.Errorf("celebration %q should have IsMovable=true", c.Name)
		}
	}

	names := make(map[string]bool)
	for _, c := range mobile {
		names[c.Name] = true
	}

	expected := []string{
		"Epifania do Senhor",
		"Batismo do Senhor",
		"Santos Pedro e Paulo, apóstolos",
		"Todos os Santos",
		"Nosso Senhor Jesus Cristo, Rei do Universo",
	}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("GetMobileDates should include %q", name)
		}
	}
}

func TestSeasonForDate(t *testing.T) {
	ly := LiturgicYearOf(2026)

	tests := []struct {
		date       Date
		wantSeason SeasonName
		wantColor  LiturgicalColor
	}{
		{NewDate(18, FEBRUARY, 2026), SeasonLent, Purple},
		{NewDate(5, APRIL, 2026), SeasonEasterSeason, White},
		{NewDate(24, MAY, 2026), SeasonEasterSeason, White},
		{NewDate(25, MAY, 2026), SeasonOrdinaryTimeII, Green},
		{NewDate(15, JUNE, 2026), SeasonOrdinaryTimeII, Green},
	}

	for _, tt := range tests {
		t.Run(tt.date.String(), func(t *testing.T) {
			season, color := seasonForDate(tt.date, ly.LiturgicSeasons)
			if season != tt.wantSeason {
				t.Errorf("season = %v, want %v", season, tt.wantSeason)
			}
			if color != tt.wantColor {
				t.Errorf("color = %v, want %v", color, tt.wantColor)
			}
		})
	}
}

func TestDateBeforeAfter(t *testing.T) {
	d1 := NewDate(1, JANUARY, 2026)
	d2 := NewDate(2, JANUARY, 2026)

	if !d1.Before(d2) {
		t.Error("Jan 1 should be before Jan 2")
	}
	if d2.Before(d1) {
		t.Error("Jan 2 should not be before Jan 1")
	}
	if !d2.After(d1) {
		t.Error("Jan 2 should be after Jan 1")
	}
	if d1.After(d2) {
		t.Error("Jan 1 should not be after Jan 2")
	}
	if d1.Before(d1) {
		t.Error("same date should not be before itself")
	}
	if d1.After(d1) {
		t.Error("same date should not be after itself")
	}
}

func TestDateRangeContains(t *testing.T) {
	r := DateRange{
		Start: NewDate(1, MARCH, 2026),
		End:   NewDate(31, MARCH, 2026),
	}

	if !r.Contains(NewDate(1, MARCH, 2026)) {
		t.Error("range should contain its start")
	}
	if !r.Contains(NewDate(31, MARCH, 2026)) {
		t.Error("range should contain its end")
	}
	if !r.Contains(NewDate(15, MARCH, 2026)) {
		t.Error("range should contain a middle date")
	}
	if r.Contains(NewDate(28, FEBRUARY, 2026)) {
		t.Error("range should not contain a date before start")
	}
	if r.Contains(NewDate(1, APRIL, 2026)) {
		t.Error("range should not contain a date after end")
	}
}

func TestGetCalendarDayMultipleCelebrations(t *testing.T) {
	// Find a day with multiple celebrations by scanning the full year
	entries, err := GetCalendar(2026)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	multiCount := 0
	for _, e := range entries {
		if len(e.Celebrations) > 1 {
			multiCount++
			// Verify sorting
			for i := 1; i < len(e.Celebrations); i++ {
				if e.Celebrations[i].Level < e.Celebrations[i-1].Level {
					t.Errorf("date %v: not sorted by precedence", e.Date)
				}
			}
		}
	}

	t.Logf("Found %d days with multiple celebrations in 2026", multiCount)
}
