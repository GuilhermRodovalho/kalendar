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

func TestEpiphany(t *testing.T) {
	tests := []struct {
		year     int
		expected Date
	}{
		// 2026: Jan 2 is Friday -> next Sunday = Jan 4
		{2026, NewDate(4, JANUARY, 2026)},
		// 2025: Jan 2 is Thursday -> next Sunday = Jan 5
		{2025, NewDate(5, JANUARY, 2025)},
		// 2028: Jan 2 is Sunday -> Jan 2
		{2028, NewDate(2, JANUARY, 2028)},
	}

	for _, tt := range tests {
		got := epiphany(tt.year)
		if got != tt.expected {
			t.Errorf("epiphany(%d): expected %v, got %v", tt.year, tt.expected, got)
		}
	}
}

func TestSaintsPeterAndPaul(t *testing.T) {
	tests := []struct {
		year     int
		expected Date
	}{
		// 2026: June 28 is Sunday -> June 28
		{2026, NewDate(28, JUNE, 2026)},
		// 2025: June 28 is Saturday -> next Sunday = June 29
		{2025, NewDate(29, JUNE, 2025)},
		// 2027: June 28 is Monday -> next Sunday = July 4
		{2027, NewDate(4, JULY, 2027)},
	}

	for _, tt := range tests {
		got := saintsPeterAndPaul(tt.year)
		if got != tt.expected {
			t.Errorf("saintsPeterAndPaul(%d): expected %v, got %v", tt.year, tt.expected, got)
		}
	}
}

func TestAssumptionOfMary(t *testing.T) {
	tests := []struct {
		year     int
		expected Date
	}{
		// 2027: Aug 15 is Sunday -> Aug 15
		{2027, NewDate(15, AUGUST, 2027)},
		// 2026: Aug 15 is Saturday -> next Sunday = Aug 16
		{2026, NewDate(16, AUGUST, 2026)},
		// 2025: Aug 15 is Friday -> next Sunday = Aug 17
		{2025, NewDate(17, AUGUST, 2025)},
	}

	for _, tt := range tests {
		got := assumptionOfMary(tt.year)
		if got != tt.expected {
			t.Errorf("assumptionOfMary(%d): expected %v, got %v", tt.year, tt.expected, got)
		}
	}
}

func TestAllSaints(t *testing.T) {
	tests := []struct {
		year     int
		expected Date
	}{
		// 2026: Nov 1 is Sunday -> Nov 1
		{2026, NewDate(1, NOVEMBER, 2026)},
		// 2025: Nov 1 is Saturday -> next Sunday = Nov 2
		{2025, NewDate(2, NOVEMBER, 2025)},
		// 2027: Nov 1 is Monday -> next Sunday = Nov 7
		{2027, NewDate(7, NOVEMBER, 2027)},
	}

	for _, tt := range tests {
		got := allSaints(tt.year)
		if got != tt.expected {
			t.Errorf("allSaints(%d): expected %v, got %v", tt.year, tt.expected, got)
		}
	}
}

func TestChristTheKing(t *testing.T) {
	tests := []struct {
		year     int
		expected Date
	}{
		// 2026: Advent starts Nov 29 -> Christ the King = Nov 22
		{2026, NewDate(22, NOVEMBER, 2026)},
		// 2025: Advent starts Nov 30 -> Christ the King = Nov 23
		{2025, NewDate(23, NOVEMBER, 2025)},
		// 2027: Advent starts Nov 28 -> Christ the King = Nov 21
		{2027, NewDate(21, NOVEMBER, 2027)},
	}

	for _, tt := range tests {
		got := christTheKing(tt.year)
		if got != tt.expected {
			t.Errorf("christTheKing(%d): expected %v, got %v", tt.year, tt.expected, got)
		}
	}
}

func TestHolyFamily(t *testing.T) {
	tests := []struct {
		year     int
		expected Date
	}{
		// 2026: Dec 25 is Friday -> Dec 27 is Sunday
		{2026, NewDate(27, DECEMBER, 2026)},
		// 2025: Dec 25 is Thursday -> Dec 28 is Sunday
		{2025, NewDate(28, DECEMBER, 2025)},
		// 2022: Dec 25 is Sunday -> Dec 30 (no Sunday in 26-31 range since 25 is Sunday)
		{2022, NewDate(30, DECEMBER, 2022)},
	}

	for _, tt := range tests {
		got := holyFamily(tt.year)
		if got != tt.expected {
			t.Errorf("holyFamily(%d): expected %v, got %v", tt.year, tt.expected, got)
		}
	}
}

func TestMaryMotherOfTheChurch(t *testing.T) {
	// 2026: Easter = April 5, Pentecost = May 24, Monday = May 25
	easter2026 := NewDate(5, APRIL, 2026)
	got := maryMotherOfTheChurch(easter2026)
	expected := NewDate(25, MAY, 2026)
	if got != expected {
		t.Errorf("maryMotherOfTheChurch(2026): expected %v, got %v", expected, got)
	}

	// 2025: Easter = April 20, Pentecost = June 8, Monday = June 9
	easter2025 := NewDate(20, APRIL, 2025)
	got = maryMotherOfTheChurch(easter2025)
	expected = NewDate(9, JUNE, 2025)
	if got != expected {
		t.Errorf("maryMotherOfTheChurch(2025): expected %v, got %v", expected, got)
	}
}

func TestImmaculateHeartOfMary(t *testing.T) {
	// 2026: Easter = April 5 -> +69 = June 13 (Saturday)
	easter2026 := NewDate(5, APRIL, 2026)
	got := immaculateHeartOfMary(easter2026)
	expected := NewDate(13, JUNE, 2026)
	if got != expected {
		t.Errorf("immaculateHeartOfMary(2026): expected %v, got %v", expected, got)
	}
	if got.Weekday() != SATURDAY {
		t.Errorf("immaculateHeartOfMary should be a Saturday, got %v", got.Weekday())
	}
}

func TestNewMobileDatesInLiturgicYear(t *testing.T) {
	ly := LiturgicYearOf(2026)

	if ly.MobileDates.Epiphany.Date.Day() == 0 {
		t.Error("Epiphany should be set")
	}
	if ly.MobileDates.BaptismOfTheLord.Date.Day() == 0 {
		t.Error("BaptismOfTheLord should be set")
	}
	if ly.MobileDates.SaintsPeterAndPaul.Date.Day() == 0 {
		t.Error("SaintsPeterAndPaul should be set")
	}
	if ly.MobileDates.AssumptionOfMary.Date.Day() == 0 {
		t.Error("AssumptionOfMary should be set")
	}
	if ly.MobileDates.AllSaints.Date.Day() == 0 {
		t.Error("AllSaints should be set")
	}
	if ly.MobileDates.ChristTheKing.Date.Day() == 0 {
		t.Error("ChristTheKing should be set")
	}
	if ly.MobileDates.HolyFamily.Date.Day() == 0 {
		t.Error("HolyFamily should be set")
	}
	if ly.MobileDates.MaryMotherOfTheChurch.Date.Day() == 0 {
		t.Error("MaryMotherOfTheChurch should be set")
	}
	if ly.MobileDates.ImmaculateHeartOfMary.Date.Day() == 0 {
		t.Error("ImmaculateHeartOfMary should be set")
	}
}
