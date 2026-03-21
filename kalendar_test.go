package kalendar

import (
	"encoding/json"
	"testing"
)

func TestEasterByGaussGregorian(t *testing.T) {
	tests := []struct {
		input  int
		output Date
	}{
		{2026, NewDate(5, APRIL, 2026)},
		{2025, NewDate(20, APRIL, 2025)},
		{2024, NewDate(31, MARCH, 2024)},
		{2023, NewDate(9, APRIL, 2023)},
		{2022, NewDate(17, APRIL, 2022)},
		{2000, NewDate(23, APRIL, 2000)},
		{1961, NewDate(2, APRIL, 1961)},
	}

	for _, tt := range tests {
		result := EasterByGauss(tt.input, GREGORIAN)
		if result != tt.output {
			t.Errorf("Gregorian Easter(%d): expected %v, got %v", tt.input, tt.output, result)
		}
	}
}

func TestEasterByGaussJulian(t *testing.T) {
	tests := []struct {
		input  int
		output Date
	}{
		{2025, NewDate(7, APRIL, 2025)},
		{2024, NewDate(22, APRIL, 2024)},
		{2023, NewDate(3, APRIL, 2023)},
	}

	for _, tt := range tests {
		result := EasterByGauss(tt.input, JULIAN)
		if result != tt.output {
			t.Errorf("Julian Easter(%d): expected %v, got %v", tt.input, tt.output, result)
		}
	}
}

func TestLent(t *testing.T) {
	tests := []struct {
		input int
		start Date
		end   Date
	}{
		{2026, NewDate(18, FEBRUARY, 2026), NewDate(1, APRIL, 2026)},
		{2025, NewDate(5, MARCH, 2025), NewDate(16, APRIL, 2025)},
		{2024, NewDate(14, FEBRUARY, 2024), NewDate(27, MARCH, 2024)},
	}

	for _, tt := range tests {
		start, end := Lent(tt.input)
		if start != tt.start || end != tt.end {
			t.Errorf("Lent(%d): expected (%v, %v), got (%v, %v)", tt.input, tt.start, tt.end, start, end)
		}
	}
}

func TestLiturgicYearMobileDates(t *testing.T) {
	ly := LiturgicYearOf(2026)

	cases := []struct {
		name          string
		got           Feast
		expectedDate  Date
		expectedColor LiturgicalColor
	}{
		{"AshWednesday", ly.AshWednesday, NewDate(18, FEBRUARY, 2026), Purple},
		{"PalmSunday", ly.PalmSunday, NewDate(29, MARCH, 2026), Red},
		{"GoodFriday", ly.GoodFriday, NewDate(3, APRIL, 2026), Red},
		{"Easter", ly.Easter, NewDate(5, APRIL, 2026), White},
		{"AscensionOfTheLord", ly.AscensionOfTheLord, NewDate(14, MAY, 2026), White},
		{"Pentecost", ly.Pentecost, NewDate(24, MAY, 2026), Red},
		{"HolyTrinity", ly.HolyTrinity, NewDate(31, MAY, 2026), White},
		{"CorpusChristi", ly.CorpusChristi, NewDate(4, JUNE, 2026), White},
		{"FeastOfSacredHeart", ly.FeastOfSacredHeart, NewDate(12, JUNE, 2026), White},
	}

	for _, c := range cases {
		if c.got.Date != c.expectedDate {
			t.Errorf("%s date: expected %v, got %v", c.name, c.expectedDate, c.got.Date)
		}
		if c.got.Color != c.expectedColor {
			t.Errorf("%s color: expected %v, got %v", c.name, c.expectedColor, c.got.Color)
		}
	}
}

func TestLiturgicSeasons(t *testing.T) {
	ly := LiturgicYearOf(2026)

	seasons := []struct {
		name          string
		got           Season
		start         Date
		end           Date
		expectedColor LiturgicalColor
	}{
		{"Advent", ly.Advent, NewDate(30, NOVEMBER, 2025), NewDate(24, DECEMBER, 2025), Purple},
		{"Christmas", ly.Christmas, NewDate(25, DECEMBER, 2025), NewDate(11, JANUARY, 2026), White},
		{"OrdinaryTimeI", ly.OrdinaryTimeI, NewDate(12, JANUARY, 2026), NewDate(17, FEBRUARY, 2026), Green},
		{"Lent", ly.LiturgicSeasons.Lent, NewDate(18, FEBRUARY, 2026), NewDate(1, APRIL, 2026), Purple},
		{"EasterTriduum", ly.EasterTriduum, NewDate(2, APRIL, 2026), NewDate(4, APRIL, 2026), White},
		{"EasterSeason", ly.EasterSeason, NewDate(5, APRIL, 2026), NewDate(24, MAY, 2026), White},
		{"OrdinaryTimeII", ly.OrdinaryTimeII, NewDate(25, MAY, 2026), NewDate(28, NOVEMBER, 2026), Green},
	}

	for _, s := range seasons {
		if s.got.Start != s.start {
			t.Errorf("%s start: expected %v, got %v", s.name, s.start, s.got.Start)
		}
		if s.got.End != s.end {
			t.Errorf("%s end: expected %v, got %v", s.name, s.end, s.got.End)
		}
		if s.got.Color != s.expectedColor {
			t.Errorf("%s color: expected %v, got %v", s.name, s.expectedColor, s.got.Color)
		}
	}
}

func TestBaptismOfTheLord(t *testing.T) {
	tests := []struct {
		year     int
		expected Date
	}{
		// 2026: Epiphany (Jan 6) is Tuesday -> next Sunday = Jan 11
		{2026, NewDate(11, JANUARY, 2026)},
		// 2023: Epiphany (Jan 6) is Friday -> next Sunday = Jan 8
		{2023, NewDate(8, JANUARY, 2023)},
		// 2034: Epiphany (Jan 6) is Friday -> next Sunday = Jan 8
		{2034, NewDate(8, JANUARY, 2034)},
	}

	for _, tt := range tests {
		got := baptismOfTheLord(tt.year)
		if got != tt.expected {
			t.Errorf("baptismOfTheLord(%d): expected %v, got %v", tt.year, tt.expected, got)
		}
	}
}

func TestDateNextOrSame(t *testing.T) {
	// March 18, 2026 is a Wednesday
	wed := NewDate(18, MARCH, 2026)
	if wed.NextOrSame(WEDNESDAY) != wed {
		t.Errorf("NextOrSame on same weekday should return same date")
	}
	if wed.NextOrSame(THURSDAY) != NewDate(19, MARCH, 2026) {
		t.Errorf("NextOrSame(THURSDAY) from Wednesday should be next day")
	}
	if wed.NextOrSame(SUNDAY) != NewDate(22, MARCH, 2026) {
		t.Errorf("NextOrSame(SUNDAY) from Wednesday should be 4 days later")
	}
}

func TestDateNext(t *testing.T) {
	// March 18, 2026 is a Wednesday
	wed := NewDate(18, MARCH, 2026)
	// Next(WEDNESDAY) on a Wednesday should return the FOLLOWING Wednesday
	if wed.Next(WEDNESDAY) != NewDate(25, MARCH, 2026) {
		t.Errorf("Next on same weekday should advance 7 days, got %v", wed.Next(WEDNESDAY))
	}
	if wed.Next(THURSDAY) != NewDate(19, MARCH, 2026) {
		t.Errorf("Next(THURSDAY) from Wednesday should be next day")
	}
}

func TestDateMarshalJSON(t *testing.T) {
	d := NewDate(5, APRIL, 2026)
	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	if string(data) != `"2026-04-05"` {
		t.Errorf("expected \"2026-04-05\", got %s", string(data))
	}
}

func TestDateUnmarshalJSON(t *testing.T) {
	var d Date
	err := json.Unmarshal([]byte(`"2026-04-05"`), &d)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	expected := NewDate(5, APRIL, 2026)
	if d != expected {
		t.Errorf("expected %v, got %v", expected, d)
	}
}

func TestDateUnmarshalJSONInvalid(t *testing.T) {
	var d Date
	err := json.Unmarshal([]byte(`"not-a-date"`), &d)
	if err == nil {
		t.Error("expected error for invalid date format")
	}
}

func TestDateString(t *testing.T) {
	d := NewDate(5, APRIL, 2026)
	if d.String() != "2026-04-05" {
		t.Errorf("expected 2026-04-05, got %s", d.String())
	}
}

func TestDateGetters(t *testing.T) {
	d := NewDate(15, MARCH, 2026)
	if d.Day() != 15 {
		t.Errorf("Day: expected 15, got %d", d.Day())
	}
	if d.Month() != MARCH {
		t.Errorf("Month: expected MARCH, got %d", d.Month())
	}
	if d.Year() != 2026 {
		t.Errorf("Year: expected 2026, got %d", d.Year())
	}
}
